package main

import (
	"errors"
	"fmt"
	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
	"time"
)


func Get_id_with_re(re *proto.Sentence)(string,error){
	if re==nil{
		return "",errors.New("项目是空的")
	}
	return re.Map[".id"],nil
}
//获取首个结果id
func Get_firstid_with_re(re []*proto.Sentence)(*proto.Sentence,error){
	if len(re)<1{
		return nil,errors.New("项目小于1")
	}
	return re[0],nil
}
//Pppoe_add 禁用
func Pppoe_add(name,interfaces,username,password string,ros *routeros.Client)error{
	_,err:=ros.Run("/interface/pppoe-client/add",fmt.Sprintf("=name=%s",name),fmt.Sprintf("=interface=%s",interfaces),fmt.Sprintf("=user=%s",username),fmt.Sprintf("=password=%s",password),fmt.Sprintf("=use-peer-dns=yes"))
	if err!=nil{
		return err
	}
	return nil
}
//Vrrp_add
func Vrrp_add(name string,vrid int,ros *routeros.Client)error{
	_,err:=ros.Run("/interface/vrrp/add",fmt.Sprintf("=name=%s",name),fmt.Sprintf("=vrid=%d",vrid))
	if err!=nil{
		return err
	}
	return nil
}
//获取pppoe根据名称
func Pppoe_get_with_name(name string,ros *routeros.Client)(*proto.Sentence,error){
	r,err:=ros.Run("/interface/print",fmt.Sprintf("?name=%s",name))
	if err!=nil{
		return nil,err
	}
	re,err:=Get_firstid_with_re(r.Re)
	if err!=nil{
		return nil,err
	}
	return re,nil
}
//pppoe 禁用
func Pppoe_disable_with_id(id string,ros *routeros.Client)error{
	_,err:=ros.Run("/int/pppoe-client/disable",fmt.Sprintf("=.id=%s",id))
	if err!=nil{
		return err
	}
	return nil
}
//pppoe 启用
func Pppoe_enable_with_id(id string,ros *routeros.Client)error{
	_,err:=ros.Run("/int/pppoe-client/enable",fmt.Sprintf("=.id=%s",id))
	if err!=nil{
		return err
	}
	return nil
}
//Interface 禁用
func Interface_disable_with_name(name string,ros *routeros.Client)error{
	_,err:=ros.Run("/interface/disable",fmt.Sprintf("=.id=%s",name))
	if err!=nil{
		return err
	}
	return nil
}
//Interface 启用
func Interface_enable_with_name(name string,ros *routeros.Client)error{
	_,err:=ros.Run("/interface/enable",fmt.Sprintf("=.id=%s",name))
	if err!=nil{
		return err
	}
	return nil
}
//获取Interface根据名称
func Interface_get_with_name(name string,ros *routeros.Client)(*proto.Sentence,error){
	r,err:=ros.Run("/interface/print",fmt.Sprintf("?name=%s",name))
	if err!=nil{
		return nil,err
	}
	re,err:=Get_firstid_with_re(r.Re)
	if err!=nil{
		return nil,err
	}
	return re,nil
}
//获取addr根据名称
func Address_get_with_name(sinterface string,ros *routeros.Client)(*proto.Sentence,error){
	r,err:=ros.Run("/ip/address/print",fmt.Sprintf("?interface=%s",sinterface))
	if err!=nil{
		return nil,err
	}
	re,err:=Get_firstid_with_re(r.Re)
	if err!=nil{
		return nil,err
	}
	return re,nil
}
//Scheduler 添加一个任务
func Scheduler_add(name string,duration time.Duration,stask string,ros *routeros.Client)error{
	sduration:=fmt.Sprintf("%02d:%02d:%02d",int(duration.Hours())%60,int(duration.Minutes())%60,int(duration.Seconds())%60)
	_,err:=ros.Run("/system/scheduler/add",fmt.Sprintf("=name=%s",name),fmt.Sprintf("=interval=%s",sduration),fmt.Sprintf("=on-event=%s",stask))
	if err!=nil{
		return err
	}
	return nil
}
//DDNS 添加一个任务
func Rosddns_add(sinterface string,ddnspre string,ddnstoken string,duration time.Duration,ros *routeros.Client)error{
	stask:=fmt.Sprintf(`
#dnspod  PPPoE

:local pppoe "%s"

:local token "%s"

:local domain "rosddns.cn"
:local record "%s"

/ip dns cache flush
:local domainname ($record  . "." . $domain)
:log info ("domain:" . $domainname)
:local iponint ([:resolve $domainname])
:log info ("remote ip:" . $iponint)

:local ipinlocal [/ip address get [/ip address find interface=$pppoe] address]
:set ipinlocal [:pick $ipinlocal 0 ([len $ipinlocal] -3)]
:log info ("local ip:"  .  $ipinlocal)
:if ($iponint = $ipinlocal) do={
:log info ("ip cmp ok")
}
:if ($iponint != $ipinlocal) do={
:local url "http://www.hatrace.cn/d.php?token=$token&ip=$ipinlocal&domain=$domain&record=$record"
/log error $url
/tool fetch url=$url mode=http keep-result=no
:set iponint $ipinlocal
/log war "ddns update ok!"
}
`,sinterface,ddnstoken,ddnspre)
	err:=Scheduler_add(ddnspre,duration,stask,ros)
	if err!=nil{
		return err
	}
	return nil
}
//Nat_srcnat_add
func Nat_srcnat_add(srcaddress,outinterface,action string,ros *routeros.Client)error{
	_,err:=ros.Run("/ip/firewall/nat/add",fmt.Sprintf("=src-address=%s",srcaddress),fmt.Sprintf("=chain=%s",CHIAIN_srcnat),fmt.Sprintf("=out-interface=%s",outinterface),fmt.Sprintf("=action=%s",action),)
	if err!=nil{
		return err
	}
	return nil
}
//Pppoe_add 禁用
func Nat_dstnat_add(protocol,dstport int,ininterface,action,toaddress string,toports int,ros *routeros.Client)error{
	_,err:=ros.Run("/ip/firewall/nat/add",fmt.Sprintf("=chain=%s", CHIAIN_dstnat),fmt.Sprintf("=protocol=%d", protocol),fmt.Sprintf("=dst-port=%d", dstport),fmt.Sprintf("=in-interface=%s",ininterface),fmt.Sprintf("=action=%s",action),fmt.Sprintf("=to-addresses=%s",toaddress),fmt.Sprintf("=to-ports=%d",toports),)
	if err!=nil{
		return err
	}
	return nil
}