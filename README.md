# WatchTower

Watch tower is a library to check and manage objects health and perform the fix on the fly. 

## Usage

```go
done := make(chan bool)
tw := watchtower.New(true)

rds := watchtower.Fixable{
	Name: "redis object health",
	Err:  "redis object is nil",
	Healthy: func() bool {
		return rds != nil
	},
	FixFunc: func() error {
		rds = new(Redis)
		return nil
	},
}
db := watchtower.Fixable{
	Name: "database object health",
	Err:  "database nil",
	Healthy: func() bool {
	    return db != nil
	},
	FixFunc: func() error {
	    db = new(Database)
	    return nil
	},
}

tw.AddWatchObject(rds, db)
go tw.Run(done)

```