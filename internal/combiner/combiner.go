package combiner

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/timsexperiments/quick-ffmpeg-runner-thingy/internal/ffmpeg"
	"github.com/timsexperiments/quick-ffmpeg-runner-thingy/internal/log"
)

// CombinerConfig holds the configuration for combining video files.
type CombinerConfig struct {
	Path         string
	Sequence     int
	OutputFormat string
	Output       string
}

// Oerforms the conversion of the m3u8 file and combines the resulting video.
func RunCombiner(cfg CombinerConfig) error {
	var outfiles []string
	tmpDir := os.TempDir()

	for i := range cfg.Sequence {
		input := createInfileName(cfg, i)
		output := createOutfileName(tmpDir, cfg, i)
		err := ffmpeg.ConvertM3U8ToVideo(input, output)
		if err != nil {
			log.L().Error(fmt.Sprintf("failed to convert %s: %v", input, err))
			return err
		}
		log.L().Debug(fmt.Sprintf("Video file created: %v", output))
		outfiles = append(outfiles, output)
	}

	dir := filepath.Dir(cfg.Path)
	combinedFileName := fmt.Sprintf("%s/%s.%s", dir, cfg.Output, cfg.OutputFormat)
	err := ffmpeg.CombineVideos(outfiles, combinedFileName)

	for _, file := range outfiles {
		err := os.Remove(file)
		if err != nil {
			log.L().Error(fmt.Sprintf("failed to remove %s: %v", file, err))
		}
	}

	return err
}

func createInfileName(cfg CombinerConfig, sequence int) string {
	basePath := strings.TrimSuffix(cfg.Path, filepath.Ext(cfg.Path))
	return fmt.Sprintf("%s%d.m3u8", basePath, sequence)
}

func createOutfileName(tmpDir string, cfg CombinerConfig, sequence int) string {
	return path.Join(tmpDir, fmt.Sprintf("%d.%s", sequence, cfg.OutputFormat))
}
