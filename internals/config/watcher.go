package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/garv2003/go-load-balancer/internals/models"
	"log"
)

func WatchConfig(path string, cfg *Config, serverPool *models.ServerPool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				log.Println("[⚙️ Config] Change detected. Reloading...")
				if err := cfg.ReloadFromFile(path); err != nil {
					log.Println("[❌ Config] Reload failed:", err)
					continue
				}

				newCfg := cfg.SafeGet()
				serverPool.ClearServers()
				for _, s := range newCfg.Servers {
					serverPool.AddServer(s)
				}
				log.Println("[✅ Config] Reloaded and servers updated.")
			}
		case err := <-watcher.Errors:
			log.Println("[❌ Watcher] Error:", err)
		}
	}
}
