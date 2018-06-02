package based

import (
	"encoding/gob"
	"bytes"
	"log"
)

type Data_pre struct {
	Doctor_id string `json:"Doctor_id"`
	Disease string `json:"Disease"`
	Chemistry_name string `json:"Chemistry_name"`
	Amount int `json:"Amount"`
}

type Data_tran struct {
	Presciption_id string `json:"Presciption_id"`
	Medicine_name string `json:"Medicine_name"`
	Amount int `json:"Amount"`
	Ts uint16 `json:"Ts"`
	Site string `json:"Site"`
	Price float32 `json:"Price"`
}

type Presciption struct {
	Type int
	Presciption_id string
	Hospital_id string
	Patient_id string
	Ts uint16
	Data *Data_pre
	Ishandled bool
	Policy string
}

type Transaction struct {
	Type int
	Patient_id string
	Data *Data_tran
}

type presciption struct {
	Type int
	Presciption_id string
	Hospital_id string
	Patient_id string
	Ts uint16
	Data_pre []byte
	Ishandled bool
	Policy string
}

type transaction struct {
	Type int
	Patient_id string
	Data_tran []byte
}

func (b *Data_pre)serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func deserializeDatapre(d []byte) *Data_pre {
	dp := new(Data_pre)

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&dp)
	if err != nil {
		log.Panic(err)
	}

	return dp
}

func (b *Data_tran)serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func deserializeDatatran(d []byte) *Data_tran {
	dt := new(Data_tran)

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&dt)
	if err != nil {
		log.Panic(err)
	}

	return dt
}

func (b *Presciption)serialize() []byte {
	var result bytes.Buffer
	temp := new(presciption)

	temp.Data_pre = b.Data.serialize()
	temp.Hospital_id = b.Hospital_id
	temp.Patient_id = b.Patient_id
	temp.Ishandled = b.Ishandled
	temp.Ts = b.Ts
	temp.Type = b.Type
	temp.Presciption_id = b.Presciption_id
	temp.Policy = b.Policy

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(temp)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func deserializePrescription(d []byte) *Presciption {
	dp := new(Presciption)
	dptemp := new(presciption)

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&dptemp)
	if err != nil {
		log.Panic(err)
	}

	dp.Data = deserializeDatapre(dptemp.Data_pre)
	dp.Hospital_id = dptemp.Hospital_id
	dp.Patient_id = dptemp.Patient_id
	dp.Ishandled = dptemp.Ishandled
	dp.Ts = dptemp.Ts
	dp.Type = dptemp.Type
	dp.Presciption_id = dptemp.Presciption_id
	dp.Policy = dptemp.Policy

	return dp
}

func (b *Transaction)serialize() []byte {
	var result bytes.Buffer
	temp := new(transaction)

	temp.Data_tran = b.Data.serialize()
	temp.Patient_id = b.Patient_id
	temp.Type = b.Type

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(temp)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func deserializeTransaction(d []byte) *Transaction {
	dp := new(Transaction)
	dptemp := new(transaction)

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&dptemp)
	if err != nil {
		log.Panic(err)
	}

	dp.Data = deserializeDatatran(dptemp.Data_tran)
	dp.Patient_id = dptemp.Patient_id
	dp.Type = dptemp.Type

	return dp
}

