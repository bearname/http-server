# Http server

### Example usage
```go
import "github.com/bearname/http-server/pkg/server"
```

```go
log.WithFields(log.Fields{"url": conf.ServeRestAddress}).Info("starting the server")

	srv := server.StartServer(conf.ServeRestAddress, handler)

	server.WaitForKillSignal(killSignalChan)
	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Error(err)
		return
	}
```