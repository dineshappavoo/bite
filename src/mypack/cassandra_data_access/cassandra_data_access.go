package cassandra_data_access

import (
	"fmt"
	"github.com/gocql/gocql"
	"time"
	//"mypack/objects"
)

func GetCassandraConnection(address string, keyspace string) (*gocql.Session, error) {
	//Cassandra connection
	cluster := gocql.NewCluster(address)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.One
	conn, err := cluster.CreateSession()
	if err != nil {

		panic(fmt.Sprintf("Error connecting: %v", err))
		return nil, err
	}
	return conn, err
}

func TestConn() {
	conn, err := GetCassandraConnection("127.0.0.1", "gocassandra")

	if err != nil {
		fmt.Println("FATAL	", time.Now(), "Data base connection issue ")
		panic(fmt.Sprintf("Error connecting: %v", err))
	} else {
		fmt.Println("INFO	", time.Now(), "Database connected successfully ")
	}
	defer conn.Close()

	var code string
	var name string
	q := conn.Query(fmt.Sprintf("select movie_name, genre from GoCassandra.Movie where movie_id='%v'", "1"))
	q.Scan(&code, &name)
	fmt.Println("Movie Name : ", code, " Genre : ", name)

}
