//安全退出的代码
const btn_login_out=document.querySelector("#nav-btn2");
const form_login_out=document.querySelector("#form_login_out");
btn_login_out.addEventListener("click",()=>{
    form_login_out.submit();
    
});


//start_charge代码
const start_charge=document.querySelector("#start_charge");
const div_operation=document.querySelector("#div-present");
const body=document.querySelector("body");
let value=0;
let car_position=1;// 1表示处于等待区,2表示处于充电区,0代表处于充电站外.

start_charge.addEventListener("click",()=>{
    if(value!==0){
        return;
    }
    let form_charge=document.createElement("form");
    let diag=document.createElement("dialog");
    let select=document.createElement("select");
    let opt1=document.createElement("option");
    let opt2=document.createElement("option");
    let submit=document.createElement("button");
    let start_x=document.createElement("button");
    let start_input=document.createElement("input");
    let start_lab=document.createElement("label");
    let p=document.createElement("p");
    let p1=document.createElement("p");
    let lab=document.createElement("label");

    lab.textContent="选择充电模式:";
    start_input.min=1;
    start_input.max=100;
    start_lab.textContent="请输入充电量,单位千瓦时";
    p.appendChild(start_lab);
    p.appendChild(start_input);
    p1.appendChild(lab);
    start_input.type="number";
    start_input.name="elec_num";
    opt1.textContent="快充模式";
    opt2.textContent="慢充模式";
    diag.textContent="请设置你的充电方案:";
    diag.appendChild(form_charge);
    form_charge.appendChild(p);
    form_charge.appendChild(p1);
    form_charge.appendChild(select);
    
    select.name="select_mode";
    select.className="myselect";
    select.id="charge_select";
    form_charge.action="start_charge";
    form_charge.method="POST";
    start_x.textContent="x";
    start_x.id="start_x";
   
    submit.className="btn btn-primary";
    submit.textContent="确认";
    submit.id="startbutton";
    select.appendChild(opt1);
    select.appendChild(opt2);
    diag.appendChild(submit);
    diag.appendChild(start_x)
    div_operation.appendChild(diag);
    diag.focus();
    diag.show();
    value=1;

    submit.addEventListener("click",()=>{
        form_charge.submit();
        value=0;
        diag.remove();

    })
    start_x.addEventListener("click",()=>{
        value=0;
        diag.remove();

    });

});

//quarry的代码
const quarry=document.querySelector("#quarry");
quarry.addEventListener("click",()=>{
    if(value!==0){
        return;
    }
    //从服务器获取数据
    //fetch()根据数据个数创建select
    
    let div_quarry=document.createElement("div");
    let form_quarry=document.createElement("form");
    let quarry_select=document.createElement("select");
    //该用for语句创建option
    let opt1=document.createElement("option");
    let opt2=document.createElement("option");
    let submit=document.createElement("button");
    let exit_btn=document.createElement("button");

    div_operation.appendChild(div_quarry);
    div_quarry.appendChild(form_quarry);
    div_quarry.id="div-quarry";
    form_quarry.appendChild(quarry_select);
    form_quarry.appendChild(submit);
    quarry_select.appendChild(opt1);
    quarry_select.appendChild(opt2);
    div_quarry.appendChild(exit_btn);

    form_quarry.action="quarry";
    form_quarry.method="post";
    submit.textContent="确认";
    submit.className="btn btn-secondary";
    exit_btn.textContent='x';
    exit_btn.id="exit-btn";

    value=1;

    submit.addEventListener("click",()=>{
        value=0;
        form_quarry.submit();
        div_quarry.remove();
    });
    exit_btn.addEventListener("click",()=>{
        value=0;
         div_quarry.remove();
    });
    

});
//modify_quarry代码
const  modify_quarry=document.querySelector("#modify_quarry");
modify_quarry.addEventListener("click",()=>{
    if(value!==0){
        return;
    }
    let diag_modify=document.createElement("dialog");
    let form_modify=document.createElement("form");
    let options=document.createElement("select");
    let opt_fast=document.createElement("option");
    let opt_normal=document.createElement("option");
    let p=document.createElement("p");
    let p1=document.createElement("p");
    let lab=document.createElement("label");
    let elc_num=document.createElement("input");
    let submit=document.createElement("button");
    let start_x=document.createElement("button");
    let num_lab=document.createElement("lab");

    form_modify.action="modify_quarry";
    form_modify.method="post";
    if(car_position===1){
        lab.textContent="修改充电模式:";
        elc_num.min=1;
        elc_num.max=100;
        num_lab.textContent="请修改充电量,单位千瓦时";
        p.appendChild(num_lab);
        p.appendChild(elc_num);
        p1.appendChild(lab);
        elc_num.type="number";
        elc_num.name="modify_elecnum";
        opt_fast.textContent="快充模式";
        opt_normal.textContent="慢充模式";
        diag_modify.textContent="请修改你的充电方案:";
        diag_modify.appendChild(form_modify);
        form_modify.appendChild(p);
        form_modify.appendChild(p1);
        form_modify.appendChild(options);

        options.name="modify_mode";
        options.className="myselect";
        options.id="charge_select";
        form_modify.action="modify_charge";
        form_modify.method="POST";
        start_x.textContent="x";
        start_x.id="start_x";
   
        submit.className="btn btn-primary";
        submit.textContent="确认";
        submit.id="startbutton";
        options.appendChild(opt_fast);
        options.appendChild(opt_normal);
        diag_modify.appendChild(submit);
        diag_modify.appendChild(start_x)
        div_operation.appendChild(diag_modify);
        diag_modify.focus();
        diag_modify.show();
        value=1;
    }
    else if(car_position===2){
        lab.textContent="当前处于充电区,请先停止充电!!!";
        diag_modify.appendChild(lab);
        diag_modify.appendChild(start_x);
        start_x.textContent="x";
        start_x.id="modify_x";
        div_operation.appendChild(diag_modify);
        diag_modify.show();
        value=1;
    }
    else{

    }
    submit.addEventListener("click",()=>{
        form_modify.submit();
        value=0;
        diag_modify.remove();

    })
    start_x.addEventListener("click",()=>{
        value=0;
        diag_modify.remove();

    });


});
