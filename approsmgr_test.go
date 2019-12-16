package main

import (
	"fmt"
	"github.com/go-routeros/routeros"
	"testing"
)

func TestAddress_get_with_name(t *testing.T) {
	err:= func() error{
		roscli,err:=routeros.Dial("dbc1.rosddns.cn:8728","admin","")
		if err!=nil{
			return err
		}
		addr,err:=Address_get_with_name("pppoe-out1",roscli)
		if err!=nil{
			return err
		}
		fmt.Println(addr.Map["address"])

		return nil
	}()
	if err!=nil{
		fmt.Println(err)
	}
}
