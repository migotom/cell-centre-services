package helpers

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/migotom/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func projectPath() string {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(path))
}

// DBmock gives dockerized test DB handle, with purge function requied to cleanup docker container after use.
func DBmock() (db *mongo.Database, purge func() error, err error) {
	dbCreds := struct {
		user     string
		password string
		database string
	}{
		user:     "cell-centre",
		password: "cell-centre",
		database: "cell-centre",
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to test database docker: %v", err)
	}

	dockerOptions := dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "latest",
		Env: []string{
			"MONGO_INITDB_ROOT_USERNAME=" + dbCreds.user,
			"MONGO_INITDB_ROOT_PASSWORD=" + dbCreds.password,
			"MONGO_INITDB_DATABASE=" + dbCreds.database,
		},
		Mounts: []string{
			projectPath() + "/db/fixtures/:/fixtures",
			projectPath() + "/db/initdb.d/:/docker-entrypoint-initdb.d",
		},
	}

	// run container
	resource, err := pool.RunWithOptions(&dockerOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("could not start test database docker container: %v", err)
	}
	purge = func() error {
		if err = pool.Purge(resource); err != nil {
			return fmt.Errorf("could not purge resource of test database docker: %v", err)
		}
		return nil
	}

	// connect to database
	pool.MaxWait = 5 * time.Second
	var dbClient *mongo.Client

	if err := pool.Retry(func() error {
		var err error
		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@localhost:%s/admin", dbCreds.user, dbCreds.password, resource.GetPort("27017/tcp")))
		if dbClient, err = mongo.Connect(context.Background(), clientOptions); err != nil {
			return err
		}
		return dbClient.Ping(context.Background(), nil)
	}); err != nil {
		return nil, nil, fmt.Errorf("could not connect to test database docker: %s", err)
	}

	db = dbClient.Database(dbCreds.database)

	return
}
