package native

import (
	"crypto/rand"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/boltdb/bolt"
)

func put(db *bolt.DB, key, val []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("MyBuket"))

		// Persist bytes to users bucket.
		return b.Put(key, val)
	})
}

func Benchmark_CreatePerTx(b *testing.B) {
	_ = os.Remove("test-data.db")
	db, err := bolt.Open("test-data.db", 0600, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		b.Fatal(err)
	}
	_, err = tx.CreateBucket([]byte("MyBuket"))
	if err != nil {
		b.Error(err)
	}
	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		b.Error(err)
		tx.Rollback()
	}
	data := []byte(`aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxaaloxa aloxaaloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa`)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err = put(db, generateID(), data); err != nil {
				b.Error(err)
			}
		}
	})
}

func generateID() []byte {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return b
}

func Benchmark_PutDataByBatch(b *testing.B) {
	_ = os.Remove("test-data.db")
	db, err := bolt.Open("test-data.db", 0600, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		b.Fatal(err)
	}
	_, err = tx.CreateBucket([]byte("MyBuket"))
	if err != nil {
		b.Error(err)
	}
	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		b.Error(err)
		tx.Rollback()
	}

	data := []byte(`aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxaaloxa aloxaaloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa`)

	txw, err := db.Begin(true)
	if err != nil {
		b.Error(err)
	}
	var errW error

	b.ResetTimer()

	// concurrecy bdon worket by default
	//b.RunParallel(func(pb *testing.PB) {
	//	for pb.Next() {
	//		if errW = txw.Bucket([]byte(`MyBuket`)).Put(generateID(), data); errW != nil {
	//			b.Error(errW)
	//		}
	//	}
	//})

	bucket := txw.Bucket([]byte(`MyBuket`))
	for i := 0; i < b.N; i++ {
		if errW = bucket.Put(generateID(), data); errW != nil {
			b.Error(errW)
		}
	}

	if errW != nil {
		txw.Rollback()
	} else {
		txw.Commit()
	}
}

func Benchmark_PutDataConcurrentlyByBatch(b *testing.B) {
	_ = os.Remove("test-data.db")
	db, err := bolt.Open("test-data.db", 0600, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		b.Fatal(err)
	}

	_, err = tx.CreateBucket([]byte("MyBuket"))
	if err != nil {
		b.Error(err)
	}
	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		b.Error(err)
		tx.Rollback()
	}

	data := []byte(`aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxaaloxa aloxaaloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa`)

	txw, err := db.Begin(true)
	if err != nil {
		b.Error(err)
	}

	putCh := make(chan []byte)
	done := make(chan struct{})
	go func() {
		var errW error
		bucket := txw.Bucket([]byte(`MyBuket`))
		for key := range putCh {
			if errW := bucket.Put(key, data); errW != nil {
				b.Error(errW)
			}
		}
		if errW != nil {
			txw.Rollback()
		} else {
			txw.Commit()
		}
		close(done)
	}()

	// test for now start go routine
	putCh <- generateID()

	putCh <- generateID()
	putCh <- generateID()
	putCh <- generateID()

	b.ResetTimer()

	// concurrecy don't worket by default
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			putCh <- generateID()
		}
	})

	close(putCh)
	<-done
	db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte("MyBuket"))
		b.Logf("count keys = %d", buck.Stats().KeyN)
		return nil
	})
}

func putter(b *testing.B, db *bolt.DB) (chan []byte, chan struct{}) {
	data := []byte(`aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxaaloxa aloxaaloxa aloxa aloxa aloxa aloxa aloxa aloxa aloxa`)

	txw, err := db.Begin(true)
	if err != nil {
		b.Error(err)
	}

	putCh := make(chan []byte)
	done := make(chan struct{})
	go func() {
		var errW error
		bucket := txw.Bucket([]byte(`MyBuket`))
		for key := range putCh {
			if errW := bucket.Put(key, data); errW != nil {
				b.Error(errW)
			}
		}
		if errW != nil {
			txw.Rollback()
		} else {
			txw.Commit()
		}
		close(done)
	}()
	return putCh, done
}

// просто получение по ключу
// в итераторе еще быстрее
func Benchmark_GetData(b *testing.B) {
	_ = os.Remove("test-data-reader.db")
	db, err := bolt.Open("test-data-reader.db", 0600, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		b.Fatal(err)
	}
	_, err = tx.CreateBucket([]byte("MyBuket"))
	if err != nil {
		b.Error(err)
	}
	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		b.Error(err)
		tx.Rollback()
	}

	if b.N > 10000 {
		wg := sync.WaitGroup{}
		for j := 0; j < b.N/10000; j++ {
			wg.Add(1)
			go func(iter int) {
				defer wg.Done()
				in, done := putter(b, db)
				start := iter * 10000
				for i := start; i < start+10000; i++ {
					in <- []byte(fmt.Sprintf("test-data-key-%d", i))
				}
				close(in)
				<-done
			}(j)
		}
		wg.Wait()
	} else {
		in, done := putter(b, db)
		for i := 0; i < b.N; i++ {
			in <- []byte(fmt.Sprintf("test-data-key-%d", i))
		}

		close(in)
		<-done
	}

	txr, err := db.Begin(false)
	reader := txr.Bucket([]byte(`MyBuket`))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if dataFind := reader.Get([]byte(fmt.Sprintf("test-data-key-%d", i))); dataFind == nil {
			b.Errorf(`get nil data`)
		}
	}

	txr.Rollback()
}

func Benchmark_Cursor(b *testing.B) {
	_ = os.Remove("test-data-reader.db")
	db, err := bolt.Open("test-data-reader.db", 0600, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		b.Fatal(err)
	}
	_, err = tx.CreateBucket([]byte("MyBuket"))
	if err != nil {
		b.Error(err)
	}
	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		b.Error(err)
		tx.Rollback()
	}

	if b.N > 10000 {
		wg := sync.WaitGroup{}
		for j := 0; j < b.N/10000; j++ {
			wg.Add(1)
			go func(iter int) {
				defer wg.Done()
				in, done := putter(b, db)
				start := iter * 10000
				for i := start; i < start+10000; i++ {
					in <- []byte(fmt.Sprintf("test-data-key-%d", i))
				}
				close(in)
				<-done
			}(j)
		}
		wg.Wait()
	} else {
		in, done := putter(b, db)
		for i := 0; i < b.N; i++ {
			in <- []byte(fmt.Sprintf("test-data-key-%d", i))
		}

		close(in)
		<-done
	}

	txr, err := db.Begin(false)
	reader := txr.Bucket([]byte(`MyBuket`))

	curs := reader.Cursor()

	b.ResetTimer()

	curs.First()
	for i := 0; i < b.N-1; i++ {
		if _, dataFind := curs.Next(); dataFind == nil {
			b.Errorf(`get nil data`)
		}
	}

	txr.Rollback()
}
