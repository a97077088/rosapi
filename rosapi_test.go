package main

import (
	"fmt"
	"github.com/go-routeros/routeros"
	"strings"
	"testing"
)


//测试导出所有ip
func Test_exportaddress(t *testing.T) {
	err:= func() error{
		roscli,err:=routeros.Dial("dbc1.rosddns.cn:8728","admin","")
		if err!=nil{
			return err
		}
		for i:=1;i<=100;i++{
			r,err:=Address_get_with_params(Rosparam{
				"?interface":fmt.Sprintf("pppoe-out%d",i),
			},roscli)
			if err!=nil{
				return err
			}
			fmt.Println(strings.ReplaceAll(r.Map["address"],"/32","----3000----null----null"))
		}
		return nil
	}()
	if err!=nil{
		fmt.Println(err)
	}
}
//测试映射nat端口
func TestNat_add(t *testing.T) {
	err:= func() error{
		roscli,err:=routeros.Dial("dbc1.rosddns.cn:8728","admin","")
		if err!=nil{
			return err
		}
		for i:=1;i<=100;i++{
			err:=Nat_add(Rosparam{
				Chain:CHAIN_dstnat,
				Protocol:Tcp,
				Dst_port:30000,
				In_interface:fmt.Sprintf("pppoe-out%d",i),
				Action:ACTION_dstnat,
				To_addresses:fmt.Sprintf("192.168.10.%d",i),
				To_ports:10001,
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
