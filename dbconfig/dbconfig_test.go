package dbconfig

import(
	"testing"
	"database/sql"
)

func TestConnect(t *testing.T) {
    testCases := []struct {
        desc	string
        
    }{
        {
            desc: "",
            
        },
    }
    for _, tC := range testCases {
        t.Run(tC.desc, func(t *testing.T) {
            var db*sql.DB = Connect()
			if db!= nil {
				err := db.Ping()
				if err!= nil {
					t.Errorf("DB connection FAILED: %v", err)
                }
                defer db.Close()
            }
			else {
				t.Errorf("DB connection FAILED")
            }
		})
    }
}

