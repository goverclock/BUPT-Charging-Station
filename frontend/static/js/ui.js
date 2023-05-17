//To do 
//1 充电详单查询首先需要订单文件id,用于区分不同订单,建议采用日期,用户便于查找指定订单
//在用户登陆后将订单文件id发到客户端,然后用户根据订单id查找订单,服务器再将详细内容发送至客户端.
//2 排队号码查询,服务器需将排队号码发送至客户端
//3 排队车辆查询,服务器需将前车的等待数发送至客户端.


//向服务器发送数据
function send_data(part_url,object){
    server_addr="http://localhost:8080";
    url=server_addr+part_url;
    fetch(url , {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(object)
      })
        .then(response => response.json())
        .then(log_info => {
          console.log(log_info);
        })
        .catch(error => {
          console.error(error);
        });
}



//从服务器取数据
function receive_data(part_url,object){
    server_addr="http://localhost:8080";
    url=server_addr+part_url;
    json_object=fetch(url)
// fetch() 返回一个 promise。当我们从服务器收到响应时，
// 会使用该响应调用 promise 的 `then()` 处理器。
     .then((response) => {
  // 如果请求没有成功，我们的处理器会抛出错误。
     if (!response.ok) {
        throw new Error(`HTTP 错误：${response.status}`);
     }
  // 否则（如果请求成功），我们的处理器通过调用
  // response.text() 以获取文本形式的响应，
  // 并立即返回 `response.text()` 返回的 promise。
        return response.text();
     })
// 若成功调用 response.text()，会使用返回的文本来调用 `then()` 处理器，
// 然后我们将其拷贝到 `poemDisplay` 框中。
     .then((text) => poemDisplay.textContent = text)
// 捕获可能出现的任何错误，
// 并在 `poemDisplay` 框中显示一条消息。
      .catch((error) => poemDisplay.textContent = `数据获取失败:${error}`);
      object=JSON.parse(json_object);
      return object;
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

//start_charge代码
let start_charge_url="/charge/submit";
let charge_date={
    chargeMode:"",
    chargeAmount:""
}
const start_charge = document.querySelector("#start_charge");
const div_operation = document.querySelector("#div-present");
const div1=document.querySelector("#div1");
const div_background=document.querySelector("#div-background");
const body = document.querySelector("body");
div_operation.remove();

start_charge.addEventListener("click", () => {
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
        charge_date.chargeMode=select.value;
        charge_date.chargeAmount=start_input.value;
        send_data(start_charge_url,charge_date);
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
const queue_ind = document.querySelector("#queue_ind");
queue_ind.addEventListener("click", () => {
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    //从服务器获取数据
    //fetch()根据数据个数创建select

    let div_queue_ind = document.createElement("div");
    let form_queue_ind = document.createElement("form");
    let queue_ind_select = document.createElement("select");
    //该用for语句创建option
    let opt1 = document.createElement("option");
    let opt2 = document.createElement("option");
    let submit = document.createElement("button");
    let exit_btn = document.createElement("button");

    div_operation.appendChild(div_queue_ind);
    div_queue_ind.appendChild(form_queue_ind);
    div_queue_ind.id = "div-queue_ind";
    form_queue_ind.appendChild(queue_ind_select);
    form_queue_ind.appendChild(submit);
    queue_ind_select.appendChild(opt1);
    queue_ind_select.appendChild(opt2);
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
    });
});
//modify_queue_ind代码
const modify_queue_ind_url="/charge/cancelCharge";
let modify_date={
    modifyMode:"",
    modifyAmount:""
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
        opt_fast.textContent = "快充模式";
        opt_normal.textContent = "慢充模式";
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
        modify_date.modifyMode=options.value;
        modify_date.modifyAmount=elc_num.value;
        send_data(modify_queue_ind_url,modify_date);
        div_operation.remove();
        div1.appendChild(div_background);
        diag_modify.remove();

    })
    start_x.addEventListener("click", () => {
        value = 0;
        div_operation.remove();
        div1.appendChild(div_background);
        diag_modify.remove();

    });


});

//结束充电的代码
const end_charge=document.querySelector("#stop_charge");
end_charge.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if(value!==0){
        return;
    }
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
        diag_end.appendChild(form_end);
        form_end.appendChild(submit);
        diag_end.appendChild(end_x2);
        diag_end.show();

    }
    else if(car_position===2){//充电区
        diag_end.textContent="您当前正在充电,确定要结束充电吗?";
        diag_end.appendChild(form_end);
        form_end.appendChild(submit);
        diag_end.appendChild(end_x1);
        diag_end.show();

    }
    else{
        diag_end.textContent="确定取消充电吗?";
        diag_end.appendChild(form_end);
        form_end.appendChild(submit);
        diag_end.appendChild(end_x1);
        diag_end.show();

    }
    submit.addEventListener("click",()=>{

        value=0;
        diag_end.remove();

    });
    end_x1.addEventListener("click",()=>{
        diag_end.remove();
        value=0;

    });
    end_x2.addEventListener("click",()=>{
        diag_end.remove();
        value=0;

    });


});
