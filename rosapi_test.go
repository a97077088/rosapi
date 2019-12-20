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
		for i:=1;i<=100;i++{
			err= Pppoe_set_with_params(Rosparam{
				Id:fmt.Sprintf("pppoe-out%d",i),
				Use_peer_dns:YES,
			},roscli)
			if err!=nil{
				return err
			}
		}
		return nil
	}()
	if err!=nil{
		fmt.Println(err)
	}
}
