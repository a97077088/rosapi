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
		r,err:= Interface_get_with_params(Rosparam{
			"?=running":"true",
		},roscli)
		if err!=nil{
			return err
		}
		fmt.Println(r)
		return nil
	}()
	if err!=nil{
		fmt.Println(err)
	}
}
