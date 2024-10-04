package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/timsexperiments/quick-ffmpeg-runner-thingy/internal/log"
)

// Converts an m3u8 file to the specified output format using ffmpeg.
func ConvertM3U8ToVideo(m3u8File, outputFile string) error {
	command := []string{
		"ffmpeg",
		"-protocol_whitelist",
		"file,http,https,tcp,tls",
		"-i",
		m3u8File,
		"-c",
		"copy",
		outputFile,
	}
	log.L().Debug("ffmpeg command", "command", strings.Join(command, " "))
	cmd := exec.Command(command[0], command[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.L().Debug("ffmpeg output", "output", string(output))
		return fmt.Errorf("ffmpeg error: %w", err)
	}

	return nil
}

// Combines multiple video files into a single output file using ffmpeg.
func CombineVideos(videoFiles []string, outputFile string) error {
	log.L().Debug("Combining videos", "files", videoFiles, "output", outputFile)
	concatFile, err := createConcatFile(videoFiles, "inputs.txt")
	if err != nil {
		return fmt.Errorf("error creating concat file: %w", err)
	}

	command := []string{"ffmpeg", "-f", "concat", "-safe", "0", "-i", concatFile, "-c:v", "copy", "-c:a", "aac", outputFile}
	log.L().Debug("ffmpeg command", "command", strings.Join(command, " "))
	cmd := exec.Command(command[0], command[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.L().Debug("ffmpeg output", "output", string(output))
		return fmt.Errorf("ffmpeg combine error: %v\n%s", err, output)
	}

	if err := os.Remove(concatFile); err != nil {
		log.L().Error("failed to remove concat file", "error", err)
	}

	return nil
}

func createConcatFile(videoFiles []string, concatFile string) (string, error) {
	dir := os.TempDir()
	fileName := path.Join(dir, concatFile)
	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create concat file: %w", err)
	}
	defer file.Close()

	for _, video := range videoFiles {
		_, err := file.WriteString(fmt.Sprintf("file '%s'\n", video))
		if err != nil {
			log.L().Error("failed to write to concat file", "err", err)
		}
	}

	return fileName, nil
}

// Extracts the audio from a video file using ffmpeg.
func ExtractAudio(videoOutput, outputPath string) error {
	log.L().Debug("Extracting audio", "files", videoOutput, "output", outputPath)
	command := []string{"ffmpeg", "-i", videoOutput, "-q:a", "0", "-map", "a", outputPath}
	cmd := exec.Command(command[0], command[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.L().Debug("ffmpeg output", "output", string(output))
		return fmt.Errorf("ffmpeg error: %w", err)
	}

	return nil
}
