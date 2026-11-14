package symmetric

/*

 */
import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// PasswordToKey can be used to convert a password (string) into a byte sequence using sha256
func PasswordToKey(password string, size AESKeySize) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	return h.Sum(nil)[:size]
}

// StreamEncrypt can be used to encrypt data from a stream and write output to a stream.
// Usually StreamEncrypt is run through the EncryptStdinToFile and similar functions
func StreamEncrypt(key []byte, input *bufio.Reader, output *bufio.Writer) (int64, error) {
	var written int64
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, fmt.Errorf("error getting cipher block: %w", err)
	}
	// Make a unique iv
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, fmt.Errorf("error generating iv: %w", err)
	}

	if ivWritten, err := output.Write(iv); err != nil {
		return int64(ivWritten), fmt.Errorf("error writing iv: %w", err)
	} else {
		written += int64(ivWritten)
	}

	stream := cipher.NewCTR(block, iv[:])

	reader := &cipher.StreamReader{S: stream, R: input}

	// Copy the input to the output stream, decrypting as we go.
	if streamWritten, err := io.Copy(output, reader); err != nil {
		written += streamWritten
		return written, fmt.Errorf("error encrypting data: %w", err)
	} else if written != aes.BlockSize {
		panic(fmt.Errorf("expected %d, got %d", aes.BlockSize, len(iv)))
	} else {
		written += streamWritten
		return written, nil
	}
}

func EncryptToFile(key []byte, input *bufio.Reader, outFile string) (written int64, err error) {
	out, err := os.Create(outFile)
	if err != nil {
		return 0, fmt.Errorf("failed to open output file %s: %w", outFile, err)
	}
	defer func() {
		cErr := out.Close()
		if err == nil {
			err = cErr
		}
	}()
	buffer := bufio.NewWriter(out)
	if written, err = StreamEncrypt(key, input, buffer); err != nil {
		return written, err
	}
	return written, buffer.Flush()
}

// StreamDecrypt can be used to encrypt data from a stream and write output to a stream.
// Usually StreamDecrypt is run through the DecryptFileToStdout and similar functions
func StreamDecrypt(key []byte, input *bufio.Reader, output *bufio.Writer) (int64, error) {
	var read int64
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, fmt.Errorf("error getting cipher: %w", err)
	}
	// Read iv from input
	iv := make([]byte, aes.BlockSize)
	if ivRead, err := io.ReadFull(input, iv); err != nil {
		read += int64(ivRead)
		return read, fmt.Errorf("error getting iv: %w", err)
	} else {
		read += int64(ivRead)
	}
	stream := cipher.NewCTR(block, iv[:])

	reader := &cipher.StreamReader{S: stream, R: input}
	// Copy the input to the output stream, decrypting as we go.
	if decryptedRead, err := io.Copy(output, reader); err != nil {
		read += decryptedRead
		return read, fmt.Errorf("error while decrypting data: %w", err)
	} else {
		read += decryptedRead
		return read, nil
	}
}

// Cycle decrypts inFile with oldKey and encrypts the result with newKey into outFile. The decrypted data is streamed
// through a pipe, so it is never written to disk and both sides run at the same time.
func Cycle(oldKey []byte, inFile string, newKey []byte, outFile string) error {
	pipeReader, pipeWriter := io.Pipe()
	type returnValues struct {
		count int64
		err   error
	}
	c := make(chan returnValues, 1)

	go func() {
		count, err := EncryptToFile(newKey, bufio.NewReader(pipeReader), outFile)
		// When encrypting fails, this makes the decrypting side stop writing instead of blocking forever
		_ = pipeReader.CloseWithError(err)
		c <- returnValues{count, err}
	}()

	writer := bufio.NewWriter(pipeWriter)
	count, err := DecryptFromFile(oldKey, inFile, writer)
	if err == nil {
		err = writer.Flush()
	}
	// With a nil error, the encrypting side reads EOF after the last byte
	_ = pipeWriter.CloseWithError(err)
	encryptResults := <-c
	if err != nil {
		return fmt.Errorf("error while decrypting: %w", err)
	}
	if encryptResults.err != nil {
		return fmt.Errorf("error while encrypting: %w", encryptResults.err)
	}
	if count != encryptResults.count {
		return fmt.Errorf("%d was decrypted, but %d was encrypted", count, encryptResults.count)
	}

	return nil
}

func DecryptFromFile(key []byte, inFile string, output *bufio.Writer) (read int64, err error) {
	in, err := os.Open(inFile)
	if err != nil {
		return 0, fmt.Errorf("failed to open input file %s: %w", inFile, err)
	}
	defer func() {
		cErr := in.Close()
		if err == nil {
			err = cErr
		}
	}()

	return StreamDecrypt(key, bufio.NewReader(in), output)
}
