//start_charge代码
const start_charge=document.querySelector("#start_charge");
const div_operation=document.querySelector("#div-present");
const body=document.querySelector("body");
let value=0;

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
    
    opt1.textContent="快充模式";
    opt2.textContent="慢充模式";
    diag.textContent="选择你的充电模式:";
    diag.appendChild(form_charge);
    form_charge.appendChild(select);
    select.name="mode";
    select.className="myselect";
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