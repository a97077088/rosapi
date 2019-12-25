package rosapi

import (
	"errors"
	"fmt"
	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
	"time"
)

//参数解析
func Parse_param(rosparam Rosparam)[]string{
	r:=make([]string,0)
	for k,v:=range rosparam{
		itr:=fmt.Sprintf("=%s=%v",k,v,)
		r=append(r,itr)
	}
	return r
}
//raw参数解析
func Parse_rawparam(rosparam Rosparam)[]string{
	r:=make([]string,0)
	for k,v:=range rosparam{
		cmp:=""
		if v!=nil{
			cmp="="
		}
		itr:=fmt.Sprintf("%s%s%v",k,cmp,v,)
		r=append(r,itr)
	}
	return r
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
//Pppoe_add
func Pppoe_add(params map[string]interface{},ros *routeros.Client)error{
	cmd:=[]string{"/interface/pppoe-client/add"}
	cmd=append(cmd,Parse_param(params)...)
	_,err:=ros.Run(cmd...)
	if err!=nil{
		return err
	}
	return nil
}
//Nat_add
func Nat_add(params map[string]interface{},ros *routeros.Client)error{
	cmd:=[]string{"/ip/firewall/nat/add"}
	cmd=append(cmd,Parse_param(params)...)
	_,err:=ros.Run(cmd...)
	if err!=nil{
		return err
	}
	return nil
}
//Route_add
func Route_add(params map[string]interface{},ros *routeros.Client)error{
	cmd:=[]string{"/ip/route/add"}
	cmd=append(cmd,Parse_param(params)...)
	_,err:=ros.Run(cmd...)
	if err!=nil{
		return err
	}
	return nil
}
//Vrrp_add
func Vrrp_add(params map[string]interface{},ros *routeros.Client)error{
	cmd:=[]string{"/interface/vrrp/add"}
	cmd=append(cmd,Parse_param(params)...)
	_,err:=ros.Run(cmd...)
	if err!=nil{
		return err
	}
	return nil
}
//获取Interface根据参数
func Interface_get_with_params(params map[string]interface{},ros *routeros.Client)(*proto.Sentence,error){
	cmd:=[]string{"/interface/print"}
	cmd=append(cmd,Parse_rawparam(params)...)
	r,err:=ros.Run(cmd...)
	if err!=nil{
		return nil,err
	}
	re,err:=Get_firstid_with_re(r.Re)
	if err!=nil{
		return nil,err
	}
	return re,nil
}
//获取Interface根据参数
func Interface_list_with_params(params map[string]interface{},ros *routeros.Client)([]*proto.Sentence,error){
	cmd:=[]string{"/interface/print"}
	cmd=append(cmd,Parse_rawparam(params)...)
	r,err:=ros.Run(cmd...)
	if err!=nil{
		return nil,err
	}
	return r.Re,nil
}
//获取addr根据参数
func Address_get_with_params(params map[string]interface{},ros *routeros.Client)(*proto.Sentence,error){
	cmd:=[]string{"/ip/address/print"}
	cmd=append(cmd,Parse_rawparam(params)...)
	r,err:=ros.Run(cmd...)
	if err!=nil{
		return nil,err
	}
	re,err:=Get_firstid_with_re(r.Re)
	if err!=nil{
		return nil,err
	}
	return re,nil
}
//获取addr根据参数
func Address_list_with_params(params map[string]interface{},ros *routeros.Client)([]*proto.Sentence,error){
	cmd:=[]string{"/ip/address/print"}
	cmd=append(cmd,Parse_rawparam(params)...)
	r,err:=ros.Run(cmd...)
	if err!=nil{
		return nil,err
	}
	return r.Re,nil
}
//Interface 禁用,参数numbers=接口名字
func Interface_disable_with_params(params map[string]interface{},ros *routeros.Client)error{
	cmd:=[]string{"/interface/disable"}
	cmd=append(cmd,Parse_param(params)...)
	_,err:=ros.Run(cmd...)
	if err!=nil{
		return err
	}
	return nil
}
//Interface 启用,参数numbers=接口名字
func Interface_enable_with_name(params map[string]interface{},ros *routeros.Client)error{
	cmd:=[]string{"/interface/enable"}
	cmd=append(cmd,Parse_param(params)...)
	_,err:=ros.Run(cmd...)
	if err!=nil{
		return err
	}
	return nil
}
//Interface,参数numbers=接口名字
func Interface_set_with_params(params map[string]interface{},ros *routeros.Client)error{
	cmd:=[]string{"/interface/set"}
	cmd=append(cmd,Parse_param(params)...)
	_,err:=ros.Run(cmd...)
	if err!=nil{
		return err
	}
	return nil
}
//Interface,参数numbers=接口名字
func Pppoe_set_with_params(params map[string]interface{},ros *routeros.Client)error{
	cmd:=[]string{"/interface/pppoe-client/set"}
	cmd=append(cmd,Parse_param(params)...)
	_,err:=ros.Run(cmd...)
	if err!=nil{
		return err
	}
	return nil
}

//获取id
func ReId(re *proto.Sentence)string{
	return re.Map[".id"]
}
//获取首个结果id
func Get_firstid_with_re(re []*proto.Sentence)(*proto.Sentence,error){
	if len(re)<1{
		return nil,errors.New("项目小于1")
	}
	return re[0],nil
}







