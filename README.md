# WatchTower

Watch tower is a library to check and manage objects health and perform the fix on the fly. 

## Usage

```go
done := make(chan bool)
tw := watchtower.New(true)

redisFixable := watchtower.Fixable{
	Name: "redis object health",
	Err:  "redis object is nil",
	Healthy: func() bool {
		return rds != nil
	},
	Fix: func() error {
		rds = new(Redis)
		return nil
	},
}
dbFixable := watchtower.Fixable{
	Name: "database object health",
	Err:  "database nil",
	Healthy: func() bool {
	    return db != nil
	},
	Fix: func() error {
	    db = new(Database)
	    return nil
	},
}

tw.AddWatchObject(redisFixable, dbFixable)
go tw.Run(done)

```