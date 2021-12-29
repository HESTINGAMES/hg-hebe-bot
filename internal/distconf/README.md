# distconf

Framework for configuration variables and an abbreviation of _distributed configuration_.

Consider using this framework if you:

* Want to configure your application from multiple sources
* Want access to configuration to be quick (atomic reads)
* Want to know when configuration changes, without restarting your application

The library comes with basic configuration Readers, but most advanced environments will want to read from something
like Redis, Consul, or Zookeeper as well.

# Core Concepts

Distconf is made up of an array of Readers that are in priority of configuration preference.
When users want to get a configuration value, the readers are searched for the value in order
and the first one is returned.  The returned value has a Get() function that atomically returns
the previous value and will return a new configuration value whenever one changes.  It is ok to call Get() inside
hot path code. Even if the original config value was fetched from a remote service, like Redis, the call
to Get() will be an atomic memory access each time.

# Example use case

The example I currently use is to have an array of Readers{ Environment, config file, consul (defaults), consul (production) }
This allows me to override configuration from the env or a config file during development, use consul for dynamic
configuration, have a default consul location for configuration shared across deployments, and have a specific
consul location for production, canary, or internal development configuration.

# Example for a user

As a user of distconf, a `distconf.Distconf` object will already be created and passed to you.  You can get string, int64, time.Duration, boolean, or float objects from this Distconf.  When you ask for them, you specify a default value that is used if for some reason all backings are down.  You should ask for this value once, then get the configuration value from this distconf holder in your app.  Generally, you create a config object that is Loaded in `main` and used in your program.

```
    // app.go
    type Config struct {
        QueueSize *distconf.Int
    }

    func (c *Config) Load(dconf *distconf.Distconf) {
        c.QueueSize = dconf.Int("app.queue_size", 100)
    }
  
    type App struct {
        Config Config
    }
    
    func (a *App) Start() {
        for {
            x := sqs.Load(a.Config.QueueSize.Get())
            // ...
        }
    }

    // main.go
    func main() {
        dconf := loadDistConf()
        c := app.Config{}
        c.Load(dconf)
        q := app.App {
            Config: c,
        }
        q.Start()
    }
```

# Example creating a distconf

```
    jconf := JSONConfig {}
    jconf.RefreshFile("config.json")
    dconf := Distconf {
        Readers: []Reader {
            &Env{}, &jconf,
        },
    }
    // The variable should only be fetched a single time.  Then, store qsize somewhere and call Get() to fetch the
    // current value of queue.size
    qsize := dconf.Int("queue.size", 10)
    for {
        fmt.Println(qsize.Get())
    }
```

# Readers

The included readers only depend upon the core libraries.  If you want to implement your own Reader, it should satisfy 
the following interface.

```
    // Reader can get a []byte value for a config key
    type Reader interface {
        // Get should return the given config value.  If the value does not exist, it should return nil, nil.
        Get(key string) ([]byte, error)
        Close()
    }
```
 
 
If your reader is dynamic (the values contained can change while the application is running), then it should also
implement the Dynamic interface.

```
    // A Dynamic config can change what it thinks a value is over time.
    type Dynamic interface {
    	// Watch should execute callback function whenever the key changes.  The parameter to callback should be the
    	// key's name.
    	Watch(key string, callback func(string)) error
}
```

# Debugging

You can expose an expvar of all the configuration keys and values with `Var()` and just the keys with `Keys()`.

# LICENSE

Apache 2.0

Originally forked from https://github.com/signalfx/golib/tree/master/distconf at
3b7d7c75a7219b2f7b8ca5dd119251254f6a2a06
