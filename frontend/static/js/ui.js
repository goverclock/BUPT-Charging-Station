//To do 
//1 充电详单查询首先需要订单文件id,用于区分不同订单,建议采用日期,用户便于查找指定订单
//在用户登陆后将订单文件id发到客户端,然后用户根据订单id查找订单,服务器再将详细内容发送至客户端.
//2 排队号码查询,服务器需将排队号码发送至客户端
//3 排队车辆查询,服务器需将前车的等待数发送至客户端.

server_addr="http://localhost:8080";
let user_id=localStorage.getItem('user_id');//获取本地存储的用户id
let local_username=localStorage.getItem('username');
let tokens=localStorage.getItem('tokens');
console.log(user_id);
console.log(localStorage.getItem('tokens'));
console.log(localStorage.getItem('username'));
const user_id_text=document.querySelector("#user_id");
const money=document.querySelector("#money");

const div_operation = document.querySelector("#div-present");
const div1=document.querySelector("#div1");
const div_background=document.querySelector("#div-background");
const body = document.querySelector("body");
const statue_lab=document.createElement("label");
div_background.appendChild(statue_lab);

user_id_text.textContent=local_username;

//向服务器发送数据
function send_data(part_url,object){
    url=server_addr+part_url;
    const res=fetch(url , {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization':`${tokens}`,
          'Access-Control-Allow-Headers':'Authorization',
	      'Access-Control-Expose-Headers':'Authorization'
        },
        body: JSON.stringify(object)
      });
      return res;
}


//从服务器取数据
function receive_data(part_url,object){
    url=server_addr+part_url;
    const res=fetch(url , {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization':`${tokens}`,
          'Access-Control-Allow-Headers':'Authorization',
	      'Access-Control-Expose-Headers':'Authorization'
        },
        body: JSON.stringify(object)
      });
      return res;
}
function balance(){
    let balance_data={
      user_id:-1
    }
    balance_data.user_id=parseInt(user_id);
    const res=send_data(getbalance_url,balance_data);
    res.then(response=>response.json())
      .then(all_data=>{
      if(all_data.code===200){
       money.textContent=all_data.data.balance.toFixed(2);
     }
    });
}
//详细订单中找到当前正在执行的订单并处理
function analyse_details(object){//传入的是response中的data部分数据
    let cnt=-1;
    let charge_msg={
        charge_mode:"",
        queue_number:"",
        charge_id:-1,
        request_charge_amount:0.0,
        step:-1,
    }
    for(let i=0;i<object.length;i++){
        if(parseInt(object[i].step)!==4){
            cnt=i;
            break;
        }
    }
    if(cnt>=0){
    if(object[cnt].charge_mode===1){
       charge_msg.charge_mode="快充";
    }
    else{
       charge_msg.charge_mode="慢充";
    }
    charge_msg.charge_id=object[cnt].charge_id;
    charge_msg.queue_number=object[cnt].queue_number;
    charge_msg.step=object[cnt].step;
    charge_msg.request_charge_amount=object[cnt].request_charge_amount;
}
   console.log("step"+charge_msg.step);
    return charge_msg;

}
function time_trance(timestamp){
    let date=new Date(new Date(timestamp*1000).getTime()+8*3600*1000);

  date=date.toJSON();
  date=date.substring(0,19).replace('T',' ');
      return date;


}


//处理详细订单数据
function detail_deal(i,object){
    let string ;
    let step=["申请已提交","排队","呼号","充电","结束充电"];
    let true_or_false=['否','是'];
    switch(parseInt(i)){

        case 0:  string="详单编号: "+object.num; break;
        case 1:  string="充电桩编号: "+object.charge_id; break;
        case 2:  if(object.charge_mode===1){
            string="充电模式: "+"快充"
        }
        else{
            string="充电模式: "+"慢充"

        }
            break;
        case 3:   string="用户名: "+object.username; break;
        case 4:   string="用户id: "+object.user_id;    break;
        case 5:   string="请求充电电量: "+object.request_charge_amount+"度";   break;
        case 6:   string="实际充电电量: "+object.real_charge_amount+"度";   break;
        case 7:   string="充电时长: "+object.charge_time+"分钟";   break;
        case 8:   string="充电费用: "+object.charge_fee+"元";   break;
        case 9:   string="服务费用: "+object.service_fee+"元";   break;
        case 10:  string="总费用: "+object.tot_fee+"元";   break;
        case 11:  string="当前状态: "+step[object.step];    break;
        case 12:  string="排队号码: "+object.queue_number;   break;
        case 13:  string="订单提交时间: "+time_trance(object.subtime);    break;
        case 14:  string="进入等待区时间: "+time_trance(object.inlinetime);   break;
        case 15:  string="叫号时间: "+time_trance(object.calltime);    break;
        case 16:  string="开始充电时间: "+time_trance(object.charge_start_time) ;   break;
        case 17:  string="结束充电时间: "+time_trance(object.charge_end_time);   break;
        case 18:  if(object.terminate_flag==="true")string="用户是否主动取消订单: 是"; else{string="用户是否主动取消订单: 否"}  break;
        case 19:  if(object.terminate_time!==0){string="主动结束时间: "+time_trance(object.terminate_time);} else{string="用户未主动结束订单"} break;
        case 20:  if(object.failed_flag==="true") string="是否订单失败: 是"; else{string="是否订单失败: 否"} break;
        case 21:  string="订单失败原因: "+object.failed_msg;   break;
    }
    return string;
}

function user_statue(statue_lab){
    const queue_ind_url="/charge/details";
    const queue_ind_data={
    user_id:-1
}
    queue_ind_data.user_id=parseInt(user_id);
    const response=send_data(queue_ind_url,queue_ind_data);
    response.then(response=>response.json())
    .then(all_data=>{
       const msg=analyse_details(all_data.data);
       switch(msg.step){
        case 0: statue_lab.textContent="充电申请已提交请耐心等待";break;
        case 1: statue_lab.textContent="等待区排队中请耐心等待";break;
        case 2: statue_lab.textContent="叫号中,请开始充电";break;
        case 3: statue_lab.textContent="正在充电中";break;
        case 4: statue_lab.textContent="已结束充电";break;
        default:statue_lab.textContent="没有正在进行的充电方案";
    }

    });

}



//状态变量
let value=0;//用于控制显示区只显示一个子窗口
//charge_msg.step 0表示申请已提交  1表示排队中  2表示呼号中   3表示充电中  
//1,2在等待区,   3表示在充电区

//安全退出的代码
const btn_login_out = document.querySelector("#nav-btn2");
const form_login_out = document.querySelector("#form_login_out");
btn_login_out.addEventListener("click", () => {
    form_login_out.submit();
});

///getbalance代码
user_statue(statue_lab);
const getbalance_url="/getbalance";
balance();
setInterval(balance,5000);

//获取用户状态的代码
setInterval("user_statue(statue_lab)",5000);



//money_charge 余额充值代码
const money_charge_url="/recharge";
let money_charge_data={
    recharge_amount:0.0,
    user_id:0
}
const money_charge=document.querySelector("#money_charge");
money_charge.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    value=1;

    let diag=document.createElement("dialog");
    let start_x=document.createElement("button");
    let input=document.createElement("input");
    let lab=document.createElement("label");
    let p=document.createElement("p");
    let submit=document.createElement("button");

    submit.textContent="确认";
    submit.className="btn btn-primary";
    submit.id="recharge-submit";
    start_x.id="recharge-x";
    start_x.className="btn btn-primary";
    start_x.style.background.color="red";
    start_x.style.left="288px";
    start_x.textContent="x";
    lab.textContent="请输入充值金额: ";
    input.min=1;
    input.max=1000000000;
    input.type="number";
    p.appendChild(lab);
    p.appendChild(input);
    diag.appendChild(p);
    diag.appendChild(submit);
    diag.appendChild(start_x);
    div_operation.appendChild(diag);
    diag.show();

    submit.addEventListener("click",(event)=>{
        event.preventDefault();
        money_charge_data.recharge_amount=parseFloat(input.value);
        money_charge_data.user_id=parseInt(user_id);
        const response=send_data(money_charge_url,money_charge_data);
        
        response.then(response=>response.json())
        .then(all_data=>{
        if(all_data.code===200){
                p.remove();
                diag.textContent="充值成功";
                diag.appendChild(start_x);
                start_x.style.left="270px";
                start_x.style.top="-20px";
                balance();
            }
            else{
                p.remove();
                diag.textContent="充值失败";
                diag.appendChild(start_x);
                start_x.style.left="270px";
                start_x.style.top="-20px";
            }
            

        });
    });
       


    start_x.addEventListener("click",()=>{
        value = 0;
        div_operation.remove();
        div1.appendChild(div_background);
        diag.remove();
    });
});




//start_charge 开始充电代码
const start_charge_url="/charge/startCharge";
let start_charge_data={
    user_id:-1
}
const start_charge=document.querySelector("#start_charge");
start_charge.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    value=1;
    let diag=document.createElement("dialog");
    let start_x=document.createElement("button");
    start_x.id="start_x";
    start_x.textContent="x";
    start_x.className="btn btn-primary";
    start_x.style.background.color="red";
    diag.textContent="等待服务器响应...";
    diag.appendChild(start_x);
    div_operation.appendChild(diag);
    diag.show();
    start_charge_data.user_id=parseInt(user_id);
    const response=send_data(start_charge_url,start_charge_data);
    response.then(response=>response.json())
    .then(all_data=>{
    if(all_data.code===200){
        diag.textContent="已开始充电";
        start_x.style.left="248px";
        user_statue(statue_lab);
        diag.appendChild(start_x);
    }
    else{
        diag.textContent="请求失败了";
        start_x.style.left="248px";
        console.log(all_data.code);
        diag.appendChild(start_x);
    }
    });
    start_x.addEventListener("click",()=>{
        value = 0;
        div_operation.remove();
        div1.appendChild(div_background);
        diag.remove();
    });

});


//charge_submit代码
const charge_submit_url="/charge/submit";
let charge_date={
    charge_mode:0,
    charge_amount:0.0,
    user_id:-1
}
const charge_submit = document.querySelector("#charge_submit");
div_operation.remove();

charge_submit.addEventListener("click", () => {
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    let form_charge = document.createElement("div");
    let diag = document.createElement("dialog");
    let select = document.createElement("select");
    let opt1 = document.createElement("option");
    let opt2 = document.createElement("option");
    let submit = document.createElement("button");
    let start_x = document.createElement("button");
    let start_input = document.createElement("input");
    let start_lab = document.createElement("label");
    let p = document.createElement("p");
    let p1 = document.createElement("p");
    let lab = document.createElement("label");

    lab.textContent = "选择充电模式:";
    start_input.min = 1;
    start_input.max = 100;
    start_lab.textContent = "请输入充电量,单位千瓦时";
    p.appendChild(start_lab);
    p.appendChild(start_input);
    p1.appendChild(lab);
    start_input.type = "number";
    start_input.name = "elec_num";
    opt1.textContent = "快充";
    opt2.textContent = "慢充";
    diag.textContent = "请设置你的充电方案:";
    diag.appendChild(form_charge);
    form_charge.appendChild(p);
    form_charge.appendChild(p1);
    form_charge.appendChild(select);

    select.name = "select_mode";
    select.className = "myselect";
    select.id = "charge_select";
    start_x.textContent = "x";
    start_x.id = "start_submit";
    start_x.className="btn btn-primary";
    start_x.style.left="288px";

    submit.className = "btn btn-primary";
    submit.textContent = "确认";
    submit.id = "startbutton";
    select.appendChild(opt1);
    select.appendChild(opt2);
    diag.appendChild(submit);
    diag.appendChild(start_x)
    div_operation.appendChild(diag);
    diag.focus();
    diag.show();
    value = 1;

    submit.addEventListener("click", () => {
        value = 0;
        if(select.value==="快充"){
            charge_date.charge_mode=1;
        }
        else{
            charge_date.charge_mode=0;
        }
        charge_date.charge_amount=parseFloat(start_input.value);
        charge_date.user_id=parseInt(user_id);
        const res=send_data(charge_submit_url,charge_date);
        res.then(res=>res.json())
        .then(all_data=>{
            console.log(all_data);
            if(all_data.code===200){
                diag.textContent="提交成功";
                diag.appendChild(start_x)
                start_x.style.top="-19px";
                start_x.style.left="270px";
                start_x.style.position="relative";
            }
            else{
                diag.textContent="请求失败";
                diag.appendChild(start_x);
                start_x.style.top="-19px";
                start_x.style.left="270px";
                start_x.style.position="relative";
            }
        })
        user_statue(statue_lab);

    })
    start_x.addEventListener("click", () => {
        value = 0;
        div_operation.remove();
        div1.appendChild(div_background);
        diag.remove();

    });

});

// queue_ind的代码
const queue_ind_url="/charge/details";
const queue_ind_data={
    user_id:-1
}
const queue_ind = document.querySelector("#queue_ind");

let div_queue_ind = document.createElement("div");
let form_queue_ind = document.createElement("div");
let queue_ind_select = document.createElement("select");
    
let submit = document.createElement("button");
let exit_btn = document.createElement("button");
exit_btn.className="btn btn-primary";

queue_ind.addEventListener("click", () => {
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    //从服务器获取数据

    queue_ind_data.user_id=parseInt(user_id);
    let response_data;
    const response=receive_data(queue_ind_url,queue_ind_data);
    response.then(response=>response.json())
    .then(all_data=>{
        if(all_data.code===200){
            let data=all_data.data;
            response_data=all_data;
            let datas=all_data.data.sort((a,b)=>b.charge_start_time-a.charge_start_time);
            for(i=0;i<datas.length;i++){
                //该用for语句创建option
                let opt=document.createElement("option");
                opt.textContent=datas[i].num;
                queue_ind_select.appendChild(opt);
            }

        }
        else{

        }
    });

    div_operation.appendChild(div_queue_ind);
    div_queue_ind.appendChild(form_queue_ind);
    div_queue_ind.id = "div-queue_ind";
    form_queue_ind.appendChild(queue_ind_select);
    form_queue_ind.appendChild(submit);
    div_queue_ind.appendChild(exit_btn);

    submit.textContent = "确认";
    submit.className = "btn btn-secondary";
    exit_btn.textContent = 'x';
    exit_btn.id = "exit-btn";
    exit_btn.style.left="670px";
    const detail_div=document.createElement("div");

    value = 1;

    submit.addEventListener("click", () => {
        form_queue_ind.remove();
        console.log("sub");
         if(response_data.code===200){
            let data=response_data.data;
            let index_num=-1;
            for(let i=0 ;i<data.length; i++){
                if(data[i].num===parseInt(queue_ind_select.value)){
                    index_num=i;
                    break;
                }

            }
            for(let i=0;i<22;i++){
                if(i%2===1){
                const p=document.createElement("p");
                const lab1=document.createElement("label");
                const lab2=document.createElement("label");
                lab1.style.width="400px";
                
                
                lab1.textContent=detail_deal(i-1,data[index_num]);
                lab2.textContent=detail_deal(i,data[index_num]);
                p.appendChild(lab1);
                p.appendChild(lab2);
                p.style.marginBottom='0';
                detail_div.appendChild(p);
                }
            }
            exit_btn.style.top="0px"
            div_queue_ind.appendChild(detail_div);

        }
        else{
            div_queue_ind.textContent="详细订单信息获取失败";

        }
    },{once:true});//运行过程会有重复绑定的问题


    exit_btn.addEventListener("click", () => {
        value = 0;
        var pObjs = detail_div.childNodes;
        console.log(pObjs.length - 1);
        for (var i = pObjs.length - 1; i >= 0; i--) { // 一定要倒序，正序是删不干净的，可自行尝试
            detail_div.removeChild(pObjs[i]);
        }
        console.log(pObjs.length - 1);
        var pObjs = div_queue_ind.childNodes;
        console.log(pObjs.length - 1);
        for (var i = pObjs.length - 1; i >= 0; i--) { // 一定要倒序，正序是删不干净的，可自行尝试
            div_queue_ind.removeChild(pObjs[i]);
        }
        console.log(pObjs.length - 1);
        detail_div.remove();
        div_queue_ind.remove();
        div_operation.remove();
        queue_ind_select.options.length=0;
        div1.appendChild(div_background);
         
   });
});
//modify_queue_ind代码
const modify_queue_ind_url="/charge/changeSubmit";
let modify_date={
    charge_mode:"",
    charge_amount:"",
    user_id:""
}
let userid_data={
    user_id:-1
}
const modify_queue_ind = document.querySelector("#modify_queue_ind");
modify_queue_ind.addEventListener("click", () => {
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    let diag_modify = document.createElement("dialog");
    let form_modify = document.createElement("div");
    let options = document.createElement("select");
    let opt_fast = document.createElement("option");
    let opt_normal = document.createElement("option");
    let p = document.createElement("p");
    let p1 = document.createElement("p");
    let lab = document.createElement("label");
    let elc_num = document.createElement("input");
    let submit = document.createElement("button");
    let start_x = document.createElement("button");
    let num_lab = document.createElement("lab");
    let charge_msg;

    start_x.className="btn btn-primary";

    userid_data.user_id=parseInt(user_id);
    const res=send_data(queue_ind_url,userid_data);
    res.then(res=>res.json())
    .then(all_data=>{
        if(all_data.code===200){
            charge_msg=analyse_details(all_data.data);
            console.log("charge_msg.step"+charge_msg);
            modidy_queue();
        }
        else{
            
            lab.textContent = "获取充电状态失败";
            diag_modify.appendChild(lab);
            diag_modify.appendChild(start_x);
            start_x.textContent = "x";
            start_x.id = "modify_x";
            div_operation.appendChild(diag_modify);
            diag_modify.show();
            value = 1;

        }
    });

function modidy_queue(){

    if (parseInt(charge_msg.step) !== 3 && parseInt(charge_msg.step)!==-1&&parseInt(charge_msg.step)!==2) {
        lab.textContent = "修改充电模式:";
        elc_num.min = 1;
        elc_num.max = 100;
        num_lab.textContent = "请修改充电量,单位千瓦时";
        p.appendChild(num_lab);
        p.appendChild(elc_num);
        p1.appendChild(lab);
        elc_num.type = "number";
        elc_num.name = "modify_elecnum";
        opt_fast.textContent = "快充";
        opt_normal.textContent = "慢充";
        diag_modify.textContent = "请修改你的充电方案:";
        diag_modify.appendChild(form_modify);
        form_modify.appendChild(p);
        form_modify.appendChild(p1);
        form_modify.appendChild(options);

        options.name = "modify_mode";
        options.className = "myselect";
        options.id = "charge_select";
        start_x.textContent = "x";
        start_x.id = "start_x";

        start_x.style.position="relative";
        start_x.style.left="285px";
        start_x.style.top="-174px";

        submit.className = "btn btn-primary";
        submit.textContent = "确认";
        submit.id = "startbutton";
        options.appendChild(opt_fast);
        options.appendChild(opt_normal);
        diag_modify.appendChild(submit);
        diag_modify.appendChild(start_x)
        div_operation.appendChild(diag_modify);
        diag_modify.focus();
        diag_modify.show();
        value = 1;
    }
    else if (parseInt(charge_msg.step) === 3||parseInt(charge_msg.step)===2) {
        lab.textContent = "当前处于充电区,请先停止充电!!!";
        start_x.style.left="71px";
        diag_modify.appendChild(lab);
        diag_modify.appendChild(start_x);
        start_x.textContent = "x";
        start_x.id = "modify_x";
        div_operation.appendChild(diag_modify);
        diag_modify.show();
        value = 1;
    }
    else {
        lab.textContent = "未提交充电申请或当前充电已结束!!!";
        diag_modify.appendChild(lab);
        diag_modify.appendChild(start_x);
        start_x.textContent = "x";
        start_x.style.left="40px";
        start_x.id = "modify_x";
        div_operation.appendChild(diag_modify);
        diag_modify.show();
        value = 1;

    }
}
    submit.addEventListener("click", () => {
        value = 0;
        if(options.value==="快充"){
            modify_date.charge_mode=1;
        }
        else{
            modify_date.charge_mode=0;
        }
        modify_date.charge_amount=elc_num.value;
        modify_date.user_id=parseInt(user_id);
        const res=send_data(modify_queue_ind_url,modify_date);
        res.then(res=>res.json())
        .then(all_data=>{
            if(all_data.code===200){
                diag_modify.textContent="修改成功";
                diag_modify.appendChild(start_x);
                start_x.style.top="-19px";
                start_x.style.left="270px";
                start_x.style.position="relative";
            }
            else if(all_data.code===403){
                diag_modify.textContent="余额不足";
                diag_modify.appendChild(start_x);
                start_x.style.top="-19px";
                start_x.style.left="270px";
                start_x.style.position="relative";
            }
            else{
                diag_modify.textContent="修改失败";
                diag_modify.appendChild(start_x);
                start_x.style.top="-19px";
                start_x.style.left="270px";
                start_x.style.position="relative";
            }

        });

    });
    start_x.addEventListener("click", () => {
        value = 0;
        div_operation.remove();
        queue_ind_select.remove();
        div1.appendChild(div_background);
        diag_modify.remove();

    });


});

//queue_ind_id代码
const queue_ind_id_url="/charge/getChargingMsg";
let quque_ind_id_data={
    user_id:-1
}
const queue_ind_id=document.querySelector("#queue_ind_id");
queue_ind_id.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if(value!==0){
        return;
    }
    value=1;
    let div_queue_ind_id=document.createElement("div");
    let exit_btn=document.createElement("button");
    
    let p1=document.createElement("p");
    let p2=document.createElement("p");
    let p3=document.createElement("p");
    let p4=document.createElement("p");
    let p5=document.createElement("p");
    let p6=document.createElement("p");
    
    
    div_queue_ind_id.appendChild(p1);
    div_queue_ind_id.appendChild(p2);
    div_queue_ind_id.appendChild(p3);
    div_queue_ind_id.appendChild(p4);
    div_queue_ind_id.appendChild(p5);
    div_queue_ind_id.appendChild(p6);
   
    div_queue_ind_id.appendChild(exit_btn);
    exit_btn.textContent = 'x';
    exit_btn.id = "exit-btn";
    exit_btn.className="btn btn-primary";
    exit_btn.style.bottom="240px";
    exit_btn.style.left="810px";
    quque_ind_id_data.user_id=parseInt(user_id);
    div_operation.appendChild(div_queue_ind_id);
    const res=send_data(queue_ind_url,quque_ind_id_data);
    res.then(res=>res.json())
    .then(all_data=>{
        if(all_data.code===200){
            charge_msg=analyse_details(all_data.data);
            console.log("charge_msg.step"+charge_msg.step);
            queue_ind_id_now(charge_msg);
        }

    });

 function queue_ind_id_now(charge_msg){
    const response=send_data(queue_ind_id_url,quque_ind_id_data);
    response.then(response=>response.json())
         .then(all_data=>{
            console.log("charge_msg"+charge_msg.step);
         if(all_data.code===200 && parseInt(charge_msg.step)!==-1 && parseInt(charge_msg.step)!==4){
            p1.textContent="排队号码: "+charge_msg.queue_number;
            p1.style.color="white";
            p2.textContent="正在等待的前车数量: "+all_data.data["waiting_count"];
            p2.style.color="white";
            p3.textContent="充电模式: "+charge_msg.charge_mode;
            p3.style.color="white";
            p4.textContent="本次请求电量: "+charge_msg.request_charge_amount+"度";
            p4.style.color="white";
            p5.textContent="充电桩编号: "+charge_msg.charge_id;
            p5.style.color="white";
            switch(parseInt(charge_msg.step)){
                case 0: p6.textContent="当前充电状态: "+"已提交充电申请";break;
                case 1: p6.textContent="当前充电状态: "+"正在排队中";break;
                case 2: p6.textContent="当前充电状态: "+"正在呼号中";break;
                case 3: p6.textContent="当前充电状态: "+"正在充电中";break;
            }
            p6.style.color="white";
        }
        else{
            exit_btn.style.bottom="40px";
            p1.textContent="当前还未提交充电方案或充电已结束";
            p1.style.color="white";

        }
    });
 }
    exit_btn.addEventListener("click",()=>{
        value = 0;
        div_operation.remove();
        div1.appendChild(div_background);
        div_queue_ind_id.remove();
    })


});
//取消充电的代码;
const cancel_charge_url="/charge/cancelCharge";
let cancel_charge_data={
    user_id:-1
}
const cancel_charge=document.querySelector("#cancel_charge");
cancel_charge.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if(value!==0){
        return;
    }
    value=1;
    let diag_cancel=document.createElement("dialog");
    let submit=document.createElement("button");
    let end_x1=document.createElement("button");
    end_x1.textContent="x";
    end_x1.id="cancel-x";
    end_x1.className="btn btn-primary";
    end_x1.style.top="-48px";
    end_x1.style.left="307px";
    submit.textContent="确认";
    submit.className="btn btn-primary";
    submit.id="cancel_btn";
    submit.name="cancel_charge";
    let charge_msg;

    cancel_charge_data.user_id=parseInt(user_id);
    const res=send_data(queue_ind_url,cancel_charge_data);
    res.then(res=>res.json())
    .then(all_data=>{
        if(all_data.code===200){
            charge_msg=analyse_details(all_data.data);
            cancel_charge();
        }
        else{
            
            div_operation.appendChild(diag_cancel);
            diag_cancel.textContent="获取充电状态失败";
            diag_cancel.appendChild(end_x1);
            diag_cancel.show();

        }
    });

function cancel_charge(){
    if(parseInt(charge_msg.step)!==3 &&parseInt(charge_msg.step)!==-1){

      div_operation.appendChild(diag_cancel);
      diag_cancel.textContent="确认要取消充电吗?";
      end_x1.style.left="125px";
      end_x1.style.top="-20px";
      

      diag_cancel.appendChild(submit);
      diag_cancel.appendChild(end_x1);
      diag_cancel.show();
    }
    else{
      div_operation.appendChild(diag_cancel);
      diag_cancel.textContent="未提交充电方案或已结束充电或处于充电状态";
      diag_cancel.appendChild(submit);
      diag_cancel.appendChild(end_x1);
      submit.remove();
      diag_cancel.show();
    }
}

    submit.addEventListener("click",(event)=>{
        event.preventDefault();
    cancel_charge_data.user_id=parseInt(user_id);
    const  response=send_data(cancel_charge_url,cancel_charge_data);
        response.then(response=>response.json())
         .then(all_data=>{
         if(all_data.code===200){
            diag_cancel.textContent="已取消充电";
            user_statue(statue_lab);
            submit.remove();
            diag_cancel.appendChild(end_x1);
            end_x1.style.left="250px";
            end_x1.style.top="-20px";
         }
         else{
            diag_cancel.textContent="服务器繁忙,请重新操作";
            submit.remove();
            diag_cancel.appendChild(end_x1);
         }
        });

    });
    end_x1.addEventListener("click",()=>{
        value=0;
        div_operation.remove();
        div1.appendChild(div_background);
        diag_cancel.remove();
    });

});

//结束充电的代码
const end_charge_url="/charge/endCharge";
let end_charge_data={
    user_id:-1
}
const end_charge=document.querySelector("#stop_charge");
end_charge.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if(value!==0){
        return;
    }
    value=1;
    let diag_end=document.createElement("dialog");
    let form_end=document.createElement("div");
    let submit=document.createElement("button");
    let end_x1=document.createElement("button");

    end_x1.textContent="x";
    end_x1.id="end-x";
    end_x1.className="btn btn-primary";
    submit.textContent="确认";
    submit.className="btn btn-primary";
    submit.id="end_btn";
    submit.name="exit_charge";
    let charge_msg;
   
    div_operation.appendChild(diag_end);

    end_charge_data.user_id=parseInt(user_id);
    const res=send_data(queue_ind_url,end_charge_data);
    res.then(res=>res.json())
    .then(all_data=>{
        if(all_data.code===200){
            charge_msg=analyse_details(all_data.data);
            end_charge();
            
        }
        else{

            diag_end.textContent="获取充电状态失败";
            end_charge_data.user_id=parseInt(user_id);
            diag_end.appendChild(form_end);
            diag_end.appendChild(end_x1);
            diag_end.show();

        }
    });
     

function end_charge(){
    if(parseInt(charge_msg.step)===3){
     diag_end.textContent="确定结束充电吗?";
     user_statue(statue_lab);
     diag_end.appendChild(form_end);
     form_end.appendChild(submit);
     diag_end.appendChild(end_x1);
     diag_end.show();
    }
    else{
        end_x1.style.top="-50px"
        diag_end.textContent="当前未处于充电状态";
        end_charge_data.user_id=parseInt(user_id);
        diag_end.appendChild(form_end);
        diag_end.appendChild(end_x1);
        diag_end.show();
    }
}

    
    submit.addEventListener("click",()=>{
        value=0;
        end_charge_data.user_id=parseInt(user_id);
        send_data(end_charge_url,end_charge_data);
        div_operation.remove();
        div1.appendChild(div_background);
        diag_end.remove();

    });
    end_x1.addEventListener("click",()=>{
        diag_end.remove();
        div_operation.remove();
        div1.appendChild(div_background);
        value=0;

    });

});
