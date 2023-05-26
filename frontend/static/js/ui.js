//To do 
//1 充电详单查询首先需要订单文件id,用于区分不同订单,建议采用日期,用户便于查找指定订单
//在用户登陆后将订单文件id发到客户端,然后用户根据订单id查找订单,服务器再将详细内容发送至客户端.
//2 排队号码查询,服务器需将排队号码发送至客户端
//3 排队车辆查询,服务器需将前车的等待数发送至客户端.
console.log(localStorage.getItem('user_id'));
console.log(localStorage.getItem('tokens'));

server_addr="http://localhost:8080";
let user_id=localStorage.getItem('user_id');//获取本地存储的用户id
let tokens=localStorage.getItem('tokens');
const user_id_text=document.querySelector("#user_id");
const money=document.querySelector("#money");

user_id_text.textContent=user_id;

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
function receive_data(part_url){
    url=server_addr+part_url;
    const response=fetch(url);
    
    return response;
}



//状态变量
let value=0;//用于控制显示区只显示一个子窗口
let car_position=1;// 1表示处于等待区,2表示处于充电区,0代表处于充电站外.

//安全退出的代码
const btn_login_out = document.querySelector("#nav-btn2");
const form_login_out = document.querySelector("#form_login_out");
btn_login_out.addEventListener("click", () => {
    form_login_out.submit();
});
//money_charge 余额充值代码
const money_charge_url="/user/recharge";
let money_charge_data={
    recharge_amount:"",
    username:""
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
    start_x.textContent="x";
    lab.textContent="请输入充值金额: ";
    input.type="number";
    input.min=1;
    input.max=1000000;
    p.appendChild(lab);
    p.appendChild(input);
    diag.appendChild(p);
    diag.appendChild(submit);
    diag.appendChild(start_x);
    div_operation.appendChild(diag);
    diag.show();

    submit.addEventListener("click",()=>{
        money_charge_data.recharge_amount=input.value;
        money_charge_data.username=user_id;
        const response=send_data(money_charge_url,money_charge_data);
        response.then(response=>response.json())
        .then(all_data=>{
        if(all_data.code===200){
         diag.textContent="充值成功";
         money.textContent=money.value+ money_charge_data.recharge_amount;
       }
       else{
        diag.textContent="充值失败,请重试";
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
    username:""
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
    diag.textContent="等待服务器响应...";
    diag.appendChild(start_x);
    div_operation.appendChild(diag);
    diag.show();
    const response=send_data(start_charge_url,start_charge_data);
    response.then(response=>response.json())
    .then(all_data=>{
    if(all_data.code===200){
        diag.textContent="已开始充电";
    }
    else{
        diag.textContent="服务器忙,请重试";
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
    chargeMode:"",
    chargeAmount:"",
    username:""
}
const charge_submit = document.querySelector("#charge_submit");
const div_operation = document.querySelector("#div-present");
const div1=document.querySelector("#div1");
const div_background=document.querySelector("#div-background");
const body = document.querySelector("body");
div_operation.remove();

charge_submit.addEventListener("click", () => {
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    let form_charge = document.createElement("form");
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
    form_charge.action = "start_charge";
    form_charge.method = "POST";
    start_x.textContent = "x";
    start_x.id = "start_x";

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
            charge_date.chargeMode=1;
        }
        else{
            charge_date.chargeMode=0;
        }
        charge_date.chargeAmount=start_input.value;
        charge_date.username=user_id;
        send_data(charge_submit_url,charge_date);
        div_operation.remove();
        div1.appendChild(div_background);
        diag.remove();

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
const queue_ind = document.querySelector("#queue_ind");

let div_queue_ind = document.createElement("div");
let form_queue_ind = document.createElement("form");
let queue_ind_select = document.createElement("select");
    
let submit = document.createElement("button");
let exit_btn = document.createElement("button");

queue_ind.addEventListener("click", () => {
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    //从服务器获取数据
    const response=receive_data(queue_ind_url);
    response.then(response=>response.json())
    .then(all_data=>{
        if(all_data.code===200){
            let data=all_data.data;
            for(i=0;i<data.length;i++){
                //该用for语句创建option
                let opt=document.createElement("option");
                opt.textContent=data[i].order_id;
                select.appendChild(opt);
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

    form_queue_ind.action = "queue_ind";
    form_queue_ind.method = "post";
    submit.textContent = "确认";
    submit.className = "btn btn-secondary";
    exit_btn.textContent = 'x';
    exit_btn.id = "exit-btn";

    value = 1;

    submit.addEventListener("click", () => {
        value = 0;
        div_queue_ind.remove();
        div_operation.remove();
        div1.appendChild(div_background);
    });
    exit_btn.addEventListener("click", () => {
        value = 0;
        div_queue_ind.remove();
        div_operation.remove();
        div1.appendChild(div_background);
         response.then(response=>response.json())
         .then(all_data=>{
         if(all_data.code===200){
            let data=all_data.data[select.value];
            const p=document.createElement("p");
            p.textContent=data;
            div_background.appendChild(p);
        }
        else{

        }
    });
});
});
//modify_queue_ind代码
const modify_queue_ind_url="/charge/changeSubmit";
let modify_date={
    modifyMode:"",
    modifyAmount:"",
    username:""
}
const modify_queue_ind = document.querySelector("#modify_queue_ind");
modify_queue_ind.addEventListener("click", () => {
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    let diag_modify = document.createElement("dialog");
    let form_modify = document.createElement("form");
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

    form_modify.action = "modify_queue_ind";
    form_modify.method = "post";
    if (car_position === 1) {
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
        form_modify.action = "modify_charge";
        form_modify.method = "POST";
        start_x.textContent = "x";
        start_x.id = "start_x";

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
    else if (car_position === 2) {
        lab.textContent = "当前处于充电区,请先停止充电!!!";
        diag_modify.appendChild(lab);
        diag_modify.appendChild(start_x);
        start_x.textContent = "x";
        start_x.id = "modify_x";
        div_operation.appendChild(diag_modify);
        diag_modify.show();
        value = 1;
    }
    else {

    }
    submit.addEventListener("click", () => {
        value = 0;
        if(options.value==="快充"){
            modify_date.modifyMode=1;
        }
        else{
            modify_date.modifyMode=0;
        }
        modify_date.modifyAmount=elc_num.value;
        modify_date.username=user_id;
        send_data(modify_queue_ind_url,modify_date);
        div_operation.remove();
        div1.appendChild(div_background);
        diag_modify.remove();

    });
    start_x.addEventListener("click", () => {
        value = 0;
        div_operation.remove();
        div1.appendChild(div_background);
        diag_modify.remove();

    });


});

//queue_ind_id代码
const queue_ind_id_url="/charge/getChargingMsg";
let quque_ind_id_data={
    username:""
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
    div_operation.appendChild(div_queue_ind_id);
    div_queue_ind_id.appendChild(p1);
    div_queue_ind_id.appendChild(p2);
    div_queue_ind_id.appendChild(p3);
    div_queue_ind_id.appendChild(p4);
    div_queue_ind_id.appendChild(p5);
    div_queue_ind_id.appendChild(exit_btn);
    exit_btn.textContent = 'x';
    exit_btn.id = "exit-btn";
    const response=send_data(queue_ind_id_url,quque_ind_id_data);
    response.then(response=>response.json())
         .then(all_data=>{
         if(all_data.code===200){
            p1.textContent="排队号码: "+quque_ind_id_data.data["queue_number"];
            p2.textContent="正在等待的前车数量: "+quque_ind_id_data.data["waiting_count"];
            if(quque_ind_id_data.data["charge_mode"]===1){
                p3.textContent="充电模式: 快充";
            }
            else{
                p3.textContent="充电模式: 慢充";
            }
            p4.textContent="本次请求充电量: "+quque_ind_id_data.data["charge_amount"];
            if(quque_ind_id_data.data["charge_state"]===0){
                p5.textContent="充电状态: 未提交充电申请";
            }
            else if(quque_ind_id_data.data["charge_state"]===1){
                p5.textContent="充电状态: 在等待区队列";
            }
            else if(quque_ind_id_data.data["charge_state"]===2){
                p5.textContent="充电状态: 在充电区队列";
            }
            else{
                p5.textContent="充电状态: 正在充电";
            }

        }
    });
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
    username:""
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
    submit.textContent="确认";
    submit.className="btn btn-primary";
    submit.id="cancel_btn";
    submit.name="cancel_charge";

    div_operation.appendChild(diag_cancel);
    diag_cancel.textContent="确认要取消充电吗?";
    diag_cancel.appendChild(submit);
    diag_cancel.appendChild(end_x1);
    diag_cancel.show();
    submit.addEventListener("click",()=>{
    const  response=send_data(cancel_charge_url,cancel_charge_data);
        response.then(response=>response.json())
         .then(all_data=>{
         if(all_data.code===200){
            diag_cancel.textContent="已取消充电";
            submit.remove();
         }
         else{
            diag_cancel.textContent="服务器繁忙,请重新操作";
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
    user_name:""
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
    let form_end=document.createElement("form");
    let submit=document.createElement("button");
    let end_x1=document.createElement("button");
    let end_x2=document.createElement("button");

    form_end.action="end_charge_action";
    form_end.method="POST";
    end_x1.textContent="x";
    end_x1.id="end-x";
    end_x2.textContent="x";
    end_x2.id="end-x2";
    end_x1.className="btn btn-primary";
    end_x2.className="btn btn-primary";
    submit.textContent="确认";
    submit.className="btn btn-primary";
    submit.id="end_btn";
    submit.name="exit_charge";
   
    div_operation.appendChild(diag_end);

    if(car_position===1){//等待区
        diag_end.textContent="您当前正处于等待区,可以直接修改充电请求,确定要结束充电吗?";
        send_data(end_charge_url,end_charge_data);
        diag_end.appendChild(form_end);
        form_end.appendChild(submit);
        diag_end.appendChild(end_x2);
        diag_end.show();

    }
    else if(car_position===2){//充电区
        diag_end.textContent="您当前正在充电,确定要结束充电吗?";
        send_data(end_charge_url,end_charge_data);
        diag_end.appendChild(form_end);
        form_end.appendChild(submit);
        diag_end.appendChild(end_x1);
        diag_end.show();

    }
    else{
        diag_end.textContent="确定结束充电吗?";
        send_data(end_charge_url,end_charge_data);
        diag_end.appendChild(form_end);
        form_end.appendChild(submit);
        diag_end.appendChild(end_x1);
        diag_end.show();

    }
    submit.addEventListener("click",()=>{
        value=0;
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
    end_x2.addEventListener("click",()=>{
        diag_end.remove();
        div_operation.remove();
        div1.appendChild(div_background);
        value=0;

    });


});
