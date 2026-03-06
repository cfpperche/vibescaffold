package onboarding

import "os"

func configDir() string {
	home, _ := os.UserHomeDir()
	return home + "/.vibescaffold"
}

func HasSeen() bool {
	_, err := os.Stat(configDir() + "/onboarding_seen")
	return err == nil
}

func MarkSeen() {
	os.MkdirAll(configDir(), 0o755)
	os.WriteFile(configDir()+"/onboarding_seen", []byte("1"), 0o644)
}
