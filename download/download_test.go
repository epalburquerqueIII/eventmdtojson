package download

import "testing"

func TestDownload(t *testing.T) {
	t.Run("RepoToDisk(<REPO_URL>,<DIR>)", func(t *testing.T) {
		t.Run("download", func(t *testing.T) {
			eventsRepo := "https://github.com/nosajio/writing"
			tmpDir := "/tmp/events"
			if _, err := RepoToDisk(eventsRepo, tmpDir); err != nil {
				t.Errorf("RepoToDisk(%s, %s) caused an error: %s", eventsRepo, tmpDir, err)
			}
		})

		t.Run("deletes", func(t *testing.T) {
			eventsRepo := "https://github.com/nosajio/writing"
			tmpDir := "/tmp/events"
			if _, err := RepoToDisk(eventsRepo, tmpDir); err != nil {
				t.Errorf("RepoToDisk(%s, %s) caused an error: %s", eventsRepo, tmpDir, err)
			}
		})

	})
}
