package based

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"strconv"
	"strings"

)

func PutPrescription(a Presciption) {
	db, err := leveldb.OpenFile("./db/Prescription.db", nil)
	if err != nil {
		log.Panic(err)
	}
	db2, err := leveldb.OpenFile("./db/Mapping.db", nil)
	if err != nil {
		log.Panic(err)
	}

	aserial := a.serialize()
	err = db.Put([]byte(a.Presciption_id), aserial, nil)
	if err != nil {
		log.Panic(err)
	}

	//病人
	data, err := db2.Get([]byte(a.Patient_id), nil)
	if err != nil {
		preids := a.Presciption_id
		err = db2.Put([]byte(a.Patient_id), []byte(preids), nil)
		if err != nil {
			log.Panic(err)
		}
	} else {
		plus := "," + a.Presciption_id
		data = append(data, []byte(plus)...)
		err = db2.Put([]byte(a.Patient_id), data, nil)
		if err != nil {
			log.Panic(err)
		}
	}

	//医院
	data, err = db2.Get([]byte(a.Hospital_id), nil)
	if err != nil {
		preids := a.Presciption_id
		err = db2.Put([]byte(a.Hospital_id), []byte(preids), nil)
		if err != nil {
			log.Panic(err)
		}
	} else {
		plus := "," + a.Presciption_id
		data = append(data, []byte(plus)...)
		err = db2.Put([]byte(a.Hospital_id), data, nil)
		if err != nil {
			log.Panic(err)
		}
	}
	defer db.Close()
	defer db2.Close()
}

func PutTransaction(a Transaction) {
	db, err := leveldb.OpenFile("./db/Transaction.db", nil)
	if err != nil {
		log.Panic(err)
	}

	aserial := a.serialize()

	last, err := db.Get([]byte("last"), nil)
	if err != nil {
		//病人id链接放置id
		err = db.Put([]byte(a.Patient_id), []byte("1"), nil)
		if err != nil {
			log.Panic(err)
		}
		//放置id链接药方信息
		err = db.Put([]byte("1"), []byte(aserial), nil)
		if err != nil {
			log.Panic(err)
		}
		//last链接最后的放置id
		err = db.Put([]byte("last"), []byte("1"), nil)
		if err != nil {
			log.Panic(err)
		}
	} else {
		//last链接最后的放置id
		no, err := strconv.Atoi(string(last))
		if err != nil {
			log.Panic(err)
		}
		plus := strconv.Itoa(no + 1)
		err = db.Put([]byte("last"), []byte(plus), nil)
		if err != nil {
			log.Panic(err)
		}
		//放置id链接药方信息
		err = db.Put([]byte(plus), []byte(aserial), nil)
		if err != nil {
			log.Panic(err)
		}
		//病人id链接放置id
		data, err := db.Get([]byte(a.Patient_id), nil)
		if err != nil {
			err = db.Put([]byte(a.Patient_id), []byte(plus), nil)
			if err != nil {
				log.Panic(err)
			}
		} else {
			data = append(data, []byte(","+plus)...)
			err = db.Put([]byte(a.Patient_id), []byte(data), nil)
			if err != nil {
				log.Panic(err)
			}
		}
	}
	defer db.Close()
}

func GetPrescriptionByid(id string) []*Presciption {
	var result []*Presciption

	db, err := leveldb.OpenFile("./db/Prescription.db", nil)
	if err != nil {
		log.Panic(err)
	}
	db2, err := leveldb.OpenFile("./db/Mapping.db", nil)
	if err != nil {
		log.Panic(err)
	}

	data, err := db2.Get([]byte(id), nil)
	if err != nil {
		log.Panic(err)
	}

	pres := strings.Split(string(data), ",")
	for _, pre := range pres {
		re, err := db.Get([]byte(pre), nil)
		if err != nil {
			log.Panic(err)
		}

		result = append(result, deserializePrescription(re))
	}
	defer db.Close()
	defer db2.Close()
	return result
}

func GetTransactionByid(id string) []*Transaction {
	var result []*Transaction

	db, err := leveldb.OpenFile("./db/Transaction.db", nil)
	if err != nil {
		log.Panic(err)
	}

	data, err := db.Get([]byte(id), nil)
	if err != nil {
		log.Panic(err)
	}

	pres := strings.Split(string(data), ",")
	for _, pre := range pres {
		re, err := db.Get([]byte(pre), nil)
		//fmt.Printf("%s\n",re)
		if err != nil {
			log.Panic(err)
		}

		result = append(result, deserializeTransaction(re))
	}
	defer db.Close()
	return result
}

func UpdatePrescription(id string) {
	db, err := leveldb.OpenFile("./db/Prescription.db", nil)
	if err != nil {
		log.Panic(err)
	}
	db2, err := leveldb.OpenFile("./db/Mapping.db", nil)
	if err != nil {
		log.Panic(err)
	}

	data, err := db2.Get([]byte(id), nil)
	if err != nil {
		log.Panic(err)
	}

	temp := deserializePrescription(data)
	temp.Ishandled = true

	err = db.Put([]byte(temp.Presciption_id), temp.serialize(), nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	defer db2.Close()
}

func GetPrescriptionByattr(attr []string) []*Presciption {
	var result []*Presciption

	db, err := leveldb.OpenFile("./db/Prescription.db", nil)
	if err != nil {
		log.Panic(err)
	}

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		value := deserializePrescription(iter.Value())
		if match(attr, value.Policy) {
			result = append(result, value)
		}
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	return result
}
