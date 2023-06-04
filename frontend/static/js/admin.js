
//状态变量
let value=0;//用于控制显示区只显示一个子窗口

const server_addr=localStorage.getItem("address");
let admin_id=localStorage.getItem('admin_id');//获取本地存储的用户id
let local_adminname=localStorage.getItem('adminname');
let tokens=localStorage.getItem('tokens');
console.log(admin_id);
console.log(localStorage.getItem('tokens'));
console.log(localStorage.getItem('adminname'));
const user_id_text=document.querySelector("#admin_id");
const money=document.querySelector("#money");

user_id_text.textContent=local_adminname;


const div_operation = document.querySelector("#div-present");
const div1=document.querySelector("#div1");
const div_background=document.querySelector("#div-background");
const body = document.querySelector("body");
div_operation.remove();

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

//生成统计报表代码
const getreport_url="/chargeports/getreport";
const getreport=document.querySelector("#getreport");
let getreport_data={
    startDate:0,
    endDate:0
}
getreport.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    value=1;
    const diag=document.createElement("dialog");
    const p1=document.createElement("p");
    const p2=document.createElement("p");
    const start_time=document.createElement("input");
    const end_time=document.createElement("input");
    const submit=document.createElement("button");
    const exit=document.createElement("button");
    const lab_start=document.createElement("label");
    const lab_end=document.createElement("label");

    lab_start.textContent="输入开始时间";
    start_time.placeholder ="2023-01-01"
    lab_end.textContent="输入结束时间";
    end_time.placeholder="2023-01-01"
    submit.className="btn btn-primary";
    submit.textContent="确认";
    submit.id="submit";
    submit.style.left="150px";
    submit.style.top="30px";
    exit.className="btn btn-primary";
    exit.id="exit";
    exit.style.left="287px";
    exit.style.bottom="121px";
    exit.style.background.color="red";
    exit.textContent="x";
    
    p1.appendChild(lab_start);
    p1.appendChild(start_time);
    p2.appendChild(lab_end);
    p2.appendChild(end_time);
    diag.appendChild(p1);
    diag.appendChild(p2);
    diag.appendChild(submit);
    diag.appendChild(exit);
    div_operation.appendChild(diag);
    diag.show();

    submit.addEventListener("click",()=>{
        var start_date = new Date(start_time.value);
        var end_date = new Date(end_time.value);
        console.log(start_date);
        console.log(end_date);
        console.log(start_date.getTime());
        console.log(end_date.getTime());

        getreport_data.startDate=parseInt(start_date.getTime());
        getreport_data.endDate=parseInt(end_date.getTime())+86400;
        const response=send_data(getreport_url,getreport_data);
        response.then(response=>response.json())
        .then(all_data=>{
            console.log(all_data);
            if(all_data.code===200){
                diag.remove();
                div_operation.appendChild(exit);
                //分析统计报表
                analyse_report(all_data.data);

            }
            else{
                p1.remove();
                p2.remove();
                submit.remove();
                diag.textContent="获取请求时间段的充电桩统计报表失败";
            }

        });


    });

    exit.addEventListener("click",()=>{
        value=0;
        var pObjs = div_operation.childNodes;
        for (var i = pObjs.length - 1; i >= 0; i--) { // 一定要倒序，正序是删不干净的，可自行尝试
            div_operation.removeChild(pObjs[i]);
        }
        div_operation.remove();
        div1.appendChild(div_background);

    });
    function analyse_report(object){
        const div_report=document.createElement("div");
        if(object===null){
            div_report.textContent="查询时段无统计报表";
            div_report.style.color="white";
        }
        else{
        for(let i=0;i<object.length;i++){
            const div_msg=document.createElement("div");
            const p1=document.createElement("p");
            const p2=document.createElement("p");
            const p3=document.createElement("p");
            const lab1=document.createElement("label");
            const lab2=document.createElement("label");
            const lab3=document.createElement("label");
            const lab4=document.createElement("label");
            const lab5=document.createElement("label");
            const lab6=document.createElement("label");
            div_msg.textContent="充电桩编号: "+object[i].charge_id;
            lab1.textContent="累计充电次数: "+object[i].tot_frequency;
            lab2.textContent="累计充电时长: "+object[i].tot_charge_time+"分钟";
            lab3.textContent="累计充电电量: "+object[i].tot_charge_amount.toFixed(2)+"度";
            lab4.textContent="累计充电费用: "+object[i].tot_charge_fee.toFixed(2)+"元";
            lab5.textContent="累计服务费用: "+object[i].tot_service_fee.toFixed(2)+"元";
            lab6.textContent="累计总费用: "+object[i].tot_tot_fee.toFixed(2)+"元";
            p1.appendChild(lab1); p1.appendChild(lab2);p1.style.marginBottom="0";
            p2.appendChild(lab3); p2.appendChild(lab4);p2.style.marginBottom="0"
            p3.appendChild(lab5); p3.appendChild(lab6);
            lab1.style.width="300px"; lab2.style.width="300px"; lab3.style.width="300px";
            lab4.style.width="300px"; lab5.style.width="300px"; lab6.style.width="300px";
            div_msg.appendChild(p1); div_msg.appendChild(p2); div_msg.appendChild(p3);
            p1.style.color="white"; p2.style.color="white"; p3.style.color="white";
            div_report.appendChild(div_msg);
        }
    }
        exit.style.left="810px";
        exit.style.top="0px";
        div_report.id="scrollable-div";
        div_operation.appendChild(div_report);

    }

});


//获取所有充电桩状态的代码
const getreports_url="/chargeports/getreports";
let getreports_data={

}
const getreports=document.querySelector("#getreports");
getreports.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    value=1;

    const div_getreports=document.createElement("div");
    div_getreports.id="scrollable-div";
    const exit=document.createElement("button");

    exit.textContent="x";
    exit.className="btn btn-primary";
    exit.id="exit";

    div_getreports.appendChild(exit);
    div_operation.appendChild(div_getreports);

    const response=send_data(getreports_url,getreports_data);
    response.then(response=>response.json())
    .then(all_data=>{
        if(all_data.code===200){
            getreports_analyse(all_data.data);
        }
        else{
            div_getreports.textContent="获取所有充电桩状态失败";
        }

    });


    exit.addEventListener("click",()=>{
        value=0;
        div_getreports.remove();
        div_operation.remove();
        div1.appendChild(div_background);
    })
    function getreports_analyse(object){
        if(object===null){
            div_getreports.textContent="无充电桩";
            div_getreports.style.color="white";
            div_getreports.appendChild(exit);
        }
        else{
        for(let i=0;i<object.length;i++){
            exit.style.left="811px";
            const p1=document.createElement("p");
            const p3=document.createElement("p");
            const lab1=document.createElement("label");
            const lab2=document.createElement("label");
            const lab3=document.createElement("label");
            const lab4=document.createElement("label");
            const lab5=document.createElement("label");
            const lab6=document.createElement("label");
            lab1.style.width="200px"; lab2.style.width="200px"; lab3.style.width="200px";
            lab4.style.width="200px"; lab5.style.width="200px"; lab6.style.width="200px";
            lab1.textContent="充电桩编号: "+object[i].charge_id;
            if(object[i].charge_mode===0){
                lab2.textContent="充电模式: "+"慢充";
            }
            else{
                lab2.textContent="充电模式: "+"快充";
            }
            switch(object[i].charge_state){
                case 0: lab3.textContent="充电桩状态: "+"空闲";break;
                case 1: lab3.textContent="充电桩状态: "+"充电中";break;
                case 2: lab3.textContent="充电桩状态: "+"关闭";break;
                default: lab3.textContent="充电桩状态: "+"故障";break;
            }
            lab4.textContent="充电总量: "+object[i].tot_charge_amount.toFixed(2)+"度";
            lab5.textContent="充电总时长: "+object[i].tot_charge_time+"分钟";
            lab6.textContent="累计充电次数"+object[i].tot_frequency;
            lab1.style.color="white"; lab2.style.color="white";
            lab3.style.color="white"; lab4.style.color="white";
            lab5.style.color="white"; lab6.style.color="white";
            p1.style.marginBottom="0";
            p1.appendChild(lab1); p1.appendChild(lab2);
            p1.appendChild(lab3); p3.appendChild(lab4);
            p3.appendChild(lab5); p3.appendChild(lab6);
            div_getreports.appendChild(p1);
            div_getreports.appendChild(p3);
            div_getreports.id="scrollable-div";
            div_operation.appendChild(div_getreports);
        }
    }

    }

});

//管理员修改充电桩开关状态代码
const switch_url="/chargeports/switch";
let switch_data={
    charge_id:0
}
const switch_charge=document.querySelector("#switch");
switch_charge.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    value=1;

    const diag=document.createElement("dialog");
    const select=document.createElement("select");
    const exit=document.createElement("button");
    const submit=document.createElement("button");
    exit.textContent="x";
    exit.id="exit";
    exit.style.left="85px";
    exit.style.bottom="20px";
    exit.className="btn btn-primary";
    submit.className="btn btn-primary";
    submit.textContent="确认";
    submit.id="submit";
    submit.style.top="100px";
    submit.style.right="150px";
    diag.appendChild(select);
    diag.appendChild(exit);
    div_operation.appendChild(diag);
    diag.show();

    const response=send_data(getreports_url,getreports_data);
    response.then(response=>response.json())
    .then(all_data=>{
        if(all_data.code===200){
            console.log(all_data.data);
            for(let i=0;i<all_data.data.length;i++){
                const p=document.createElement("option");
                let state=["空闲","充电中","关闭","故障"];
                let fast_normal=["慢充","快充"];
                if(all_data.data[i].charge_state!==3){
                    p.textContent="编号: "+all_data.data[i].charge_id+" 状态: "+state[all_data.data[i].charge_state]+" 模式: "+fast_normal[all_data.data[i].charge_mode];
                    select.appendChild(p);
                }
            }
            diag.appendChild(submit);

        }
        else{
            diag.appendChild(exit);
            diag.textContent="服务器繁忙,获取充电桩状态失败";
        }

    });
    submit.addEventListener("click",()=>{
        submit.remove();
        select.remove();
        var numbers = select.value.match(/\d+/g);
        switch_data.charge_id=parseInt(numbers);
        const response=send_data(switch_url,switch_data);
        response.then(response=>response.json())
        .then(all_data=>{
            if(all_data.code===200){
                diag.textContent="充电桩状态修改成功";
                diag.appendChild(exit);
            }
            else{
                diag.textContent="充电桩状态修改失败";
                diag.appendChild(exit);

            }

        });

    });

    exit.addEventListener("click",()=>{
        value=0;
        diag.remove();
        div_operation.remove();
        div1.appendChild(div_background);

    });

});

//充电桩故障状态切换代码
const switchBroken_url="/chargeports/switchBroken";
let switchBroken_data={
    charge_id:0
}
const switchBroken_charge=document.querySelector("#switchBroken");
switchBroken_charge.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    value=1;

    const diag=document.createElement("dialog");
    const select=document.createElement("select");
    const exit=document.createElement("button");
    const submit=document.createElement("button");
    exit.textContent="x";
    exit.id="exit";
    exit.style.left="85px";
    exit.style.bottom="20px";
    exit.className="btn btn-primary";
    submit.className="btn btn-primary";
    submit.textContent="确认";
    submit.id="submit";
    submit.style.top="100px";
    submit.style.right="150px";
    diag.appendChild(select);
    diag.appendChild(exit);
    div_operation.appendChild(diag);
    diag.show();

    const response=send_data(getreports_url,getreports_data);
    response.then(response=>response.json())
    .then(all_data=>{
        if(all_data.code===200){
            console.log(all_data.data);
            for(let i=0;i<all_data.data.length;i++){
                const p=document.createElement("option");
                let state=["空闲","充电中","关闭","故障"];
                let fast_normal=["慢充","快充"];
                    p.textContent="编号: "+all_data.data[i].charge_id+" 状态: "+state[all_data.data[i].charge_state]+" 模式: "+fast_normal[all_data.data[i].charge_mode];
                    select.appendChild(p);
            }
            diag.appendChild(submit);

        }
        else{
            diag.appendChild(exit);
            diag.textContent="服务器繁忙,获取充电桩状态失败";
        }

    });
    submit.addEventListener("click",()=>{
        submit.remove();
        select.remove();
        var numbers = select.value.match(/\d+/g);
        switchBroken_data.charge_id=parseInt(numbers);
        const response=send_data(switchBroken_url,switchBroken_data);
        response.then(response=>response.json())
        .then(all_data=>{
            if(all_data.code===200){
                diag.textContent="充电桩状态修改成功";
                diag.appendChild(exit);
            }
            else{
                diag.textContent="充电桩状态修改失败";
                diag.appendChild(exit);

            }

        });

    });

    exit.addEventListener("click",()=>{
        value=0;
        diag.remove();
        div_operation.remove();
        div1.appendChild(div_background);

    });

});

//获取充电桩排队车辆信息
const waitingCars_url="/chargeports/waitingCars"
const waitingCars_data={
    charge_id:0
}
const waitingCars=document.querySelector("#waitingCars");
waitingCars.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    value=1;

    const diag=document.createElement("dialog");
    const select=document.createElement("select");
    const exit=document.createElement("button");
    const submit=document.createElement("button");
    exit.textContent="x";
    exit.id="exit";
    exit.style.left="175px";
    exit.style.bottom="20px";
    exit.className="btn btn-primary";
    submit.className="btn btn-primary";
    submit.textContent="确认";
    submit.id="submit";
    submit.style.top="100px";
    submit.style.right="50px";
    diag.appendChild(select);
    diag.appendChild(exit);
    div_operation.appendChild(diag);
    diag.show();

    const response=send_data(getreports_url,getreports_data);
    response.then(response=>response.json())
    .then(all_data=>{
        if(all_data.code===200){
            console.log(all_data.data);
            for(let i=0;i<all_data.data.length;i++){
                const p=document.createElement("option");
                let state=["空闲","充电中","关闭","故障"];
                p.textContent="编号: "+all_data.data[i].charge_id+" 状态: "+state[all_data.data[i].charge_state];
                select.appendChild(p);
            }
            diag.appendChild(submit);

        }
        else{
            diag.appendChild(exit);
            diag.textContent="服务器繁忙,获取充电桩状态失败";
        }


      });

      submit.addEventListener("click",()=>{
        submit.remove();
        select.remove();
        var numbers = select.value.match(/\d+/g);
        waitingCars_data.charge_id=parseInt(numbers);
        const response=send_data(waitingCars_url,waitingCars_data);
        response.then(response=>response.json())
        .then(all_data=>{
            console.log(all_data);
            if(all_data.code===200){
                diag.remove();
                analyse_cars(all_data.data);

            }
            else if(all_data.code===0){
                diag.textContent="当前无排队车辆";
                diag.appendChild(exit);
                exit.style.left="213px";

            }
            else{
                
                diag.textContent="获取充电桩车辆排队信息失败";
                diag.appendChild(exit);
                exit.style.left="113px";

            }

        });

    });

    exit.addEventListener("click",()=>{
        value=0;
        var pObjs = div_operation.childNodes;
        for (var i = pObjs.length - 1; i >= 0; i--) { // 一定要倒序，正序是删不干净的，可自行尝试
            div_operation.removeChild(pObjs[i]);
        }
        div_operation.remove();
        div1.appendChild(div_background);

    });
    function analyse_cars(object){
        const div_car_msg=document.createElement("div");
        div_car_msg.id="scrollable-div";
        div_car_msg.appendChild(exit);
        if(object===null){
            const p1=document.createElement("p");
            p1.textContent="当前无排队车辆";
            div_car_msg.appendChild(p1);
            p1.style.color="white";
            
        }
        else{
          for(let i=0;i<object.length;i++){
            const p1=document.createElement("p");
            const p2=document.createElement("p");
            const lab1=document.createElement("label");
            const lab2=document.createElement("label");
            const lab3=document.createElement("label");
            const lab4=document.createElement("label");
            const lab5=document.createElement("label");
            lab1.textContent="用户名: "+object[i].username;
            lab2.textContent="用户id: "+object[i].user_id;
            lab3.textContent="排队时长: "+object[i].waiting_time;
            lab4.textContent="请求充电量; "+object[i].charge_amount;
            lab5.textContent="电池容量: "+object[i].battery_capacity;
            lab1.style.width="300px"; lab2.style.width="300px";
            lab3.style.width="200px"; lab4.style.width="200px"; lab5.style.width="200px"; 
            p1.appendChild(lab1); p1.appendChild(lab2); p1.style.marginBottom="0";
            p2.appendChild(lab3); p2.appendChild(lab4); p2.appendChild(lab5); 
            p1.style.color="white"; p2.style.color="white";
            div_car_msg.appendChild(p1); div_car_msg.appendChild(p2);
        }
    }
        div_car_msg.id="scrollable-div";
        div_operation.appendChild(div_car_msg);

    }
});

//获取系统参数代码
const getsettings_url="/system/getsettings";
let getsettings_data={

}
const getsettings=document.querySelector("#getsettings");
getsettings.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    value=1;

    const div_getsettings=document.createElement("div");
    const exit=document.createElement("exit");
    div_getsettings.id="scrollable-div";
    div_getsettings.appendChild(exit);

    exit.textContent="x";
    exit.id="exit";
    exit.style.left="175px";
    exit.style.bottom="20px";
    exit.className="btn btn-primary";

    const response=send_data(getsettings_url,getsettings_data);
    response.then(response=>response.json())
    .then(all_data=>{
        if(all_data.code===200){
            exit.style.top="0px";
            exit.style.left="810px";
            const fault_schedule=["优先级调度","时间调度"];
            const schedule=["被调度车辆完成充电所需时长最短","单次调度总充电时长最短","批量调度总充电时长最短"]
            const p1=document.createElement("p");
            const p2=document.createElement("p");
            const p3=document.createElement("p");
            const p4=document.createElement("p");
            p1.textContent="等待区容量: "+all_data.data.waiting_area_size;
            p2.textContent="充电桩排队队列最大长度: "+all_data.data.charging_queue_len;
            p3.textContent="调度策略: "+schedule[all_data.data.call_schedule];
            p4.textContent="调度策略: "+fault_schedule[all_data.data.fault_schedule];
            p1.style.color="white"; p2.style.color="white"; p3.style.color="white";
            p4.style.color="white";
            div_getsettings.appendChild(p1); div_getsettings.appendChild(p2);
            div_getsettings.appendChild(p3); div_getsettings.appendChild(p4);
            div_operation.appendChild(div_getsettings);
        }
        else{
            div_getsettings.textContent="获取系统参数失败";
        }

    });

    exit.addEventListener("click",()=>{
        value=0;
        var pObjs = div_operation.childNodes;
        for (var i = pObjs.length - 1; i >= 0; i--) { // 一定要倒序，正序是删不干净的，可自行尝试
            div_operation.removeChild(pObjs[i]);
        }
        div_operation.remove();
        div1.appendChild(div_background);

    });

});

//修改系统参数代码
const  setsettings_url="/system/setsettings";
const setsettings_data={
    waiting_area_size:0,
    charging_queue_len:0,
    call_schedule:-1,
    fault_schedule:-1
}
const  setsettings=document.querySelector("#setsettings");
setsettings.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    value=1;

    const diag=document.createElement("dialog");
    const select_schedule=document.createElement("select");
    const select_fault_schedule=document.createElement("select");
    const submit=document.createElement("button");
    const exit=document.createElement("button");
    const p1=document.createElement("p");
    const input_waiting=document.createElement("input");
    const lab1=document.createElement("label");
    input_waiting.type="number";
    input_waiting.min=1;

    const lab2=document.createElement("label");
    const p2=document.createElement("p");
    const input_queue=document.createElement("input");
    const opt1=document.createElement("option");//select_schedule
    const opt2=document.createElement("option");//select_schedule
    const opt3=document.createElement("option");//select_schedule
    const opt4=document.createElement("option");//select_fault_schedule
    const opt5=document.createElement("option");//select_fault_schedule
    
    lab1.textContent="等待区容量: ";
    lab2.textContent="充电桩排队队列最大长度: ";
    opt1.textContent="被调度车辆完成充电所需时长最短";
    opt2.textContent="单次调度总充电时长最短";
    opt3.textContent="批量调度总充电时长最短";
    opt4.textContent="优先级调度";
    opt5.textContent="时间顺序调度";

    input_queue.type="number";
    input_queue.min=1;
    input_queue.max=255;



    exit.textContent="x";
    exit.id="exit";
    exit.style.left="175px";
    exit.style.bottom="20px";
    exit.className="btn btn-primary";
    submit.className="btn btn-primary";
    submit.textContent="确认";
    submit.id="submit";
    submit.style.top="100px";
    submit.style.right="50px";
    
    p1.appendChild(lab1); p1.appendChild(input_waiting);
    p2.appendChild(lab2); p2.appendChild(input_queue);
    select_schedule.appendChild(opt1); select_schedule.appendChild(opt2); select_schedule.appendChild(opt3);
    select_fault_schedule.appendChild(opt4); select_fault_schedule.appendChild(opt5);
    diag.appendChild(p1);
    diag.appendChild(p2);
    diag.appendChild(select_schedule);
    diag.appendChild(select_fault_schedule);
    diag.appendChild(submit);
    diag.appendChild(exit);
    exit.style.left="146px";
    exit.style.bottom="150px";
    submit.style.right="-10px";

    div_operation.appendChild(diag);
    diag.show();

    submit.addEventListener("click",()=>{
        const cnt1=select_index(select_schedule,select_schedule.value);
        const cnt2=select_index(select_fault_schedule,select_fault_schedule.value);
        console.log("cnt1"+cnt1);console.log("cnt2"+cnt2);
        setsettings_data.waiting_area_size=parseInt(input_waiting.value);
        setsettings_data.charging_queue_len=parseInt(input_queue.value);
        setsettings_data.call_schedule=cnt1;
        setsettings_data.fault_schedule=cnt2;
        const response=send_data(setsettings_url,setsettings_data);
        response.then(response=>response.json())
        .then(all_data=>{
            if(all_data.code===200){
                diag.textContent="系统参数修改成功";
                diag.appendChild(exit);
                exit.style.left="190px";
            }
            else{
            
                diag.textContent="系统参数修改失败";
                diag.appendChild(exit);
                exit.style.left="190px";
            }

        });


    });

    exit.addEventListener("click",()=>{
        value=0;
        var pObjs = div_operation.childNodes;
        for (var i = pObjs.length - 1; i >= 0; i--) { // 一定要倒序，正序是删不干净的，可自行尝试
            div_operation.removeChild(pObjs[i]);
        }
        div_operation.remove();
        div1.appendChild(div_background);

    });

    function select_index(select,value){
        let cnt=-1;
        for(let i=0;i<select.options.length;i++){
            if(select.options[i].value===value){
                cnt=i;
                break;
            }
        }
        return cnt;

    }


});
//添加充电桩
const addchargeport_url="/chargeports/addchargeport";
let addchargeport_data={
    charge_mode:-1
}
const  addchargeport=document.querySelector("#addchargeport");
addchargeport.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    value=1;

    const diag=document.createElement("dialog");
    const select=document.createElement("select");
    const exit=document.createElement("button");
    const submit=document.createElement("button");
    const opt1=document.createElement("option");
    const opt2=document.createElement("option");
    exit.textContent="x";
    exit.id="exit";
    exit.style.left="27px";
    exit.style.bottom="20px";
    exit.className="btn btn-primary";
    submit.className="btn btn-primary";
    submit.textContent="确认";
    submit.id="submit";
    submit.style.top="100px";
    submit.style.left="150px";
    opt1.textContent="快充";
    opt2.textContent="慢充";
    diag.textContent="请选择要添加的充电桩的类型:";
    select.appendChild(opt1); select.appendChild(opt2);
    diag.appendChild(select);
    diag.appendChild(exit);
    diag.appendChild(submit);
    div_operation.appendChild(diag);
    diag.show();
    
    submit.addEventListener("click",()=>{
      if(select.value==="快充"){
        addchargeport_data.charge_mode=1;
      }
      else{
        addchargeport_data.charge_mode=0;
      }
      const res=send_data(addchargeport_url,addchargeport_data);
      res.then(res=>res.json())
      .then(all_data=>{
        console.log(all_data);
        if(all_data.code===200){
            diag.textContent="添加成功";
            diag.appendChild(exit);
            exit.style.top="-19px";
            exit.style.left="270px";
            exit.style.position="relative";
        }
        else{
            diag.textContent="添加失败";
            diag.appendChild(exit);
            exit.style.top="-19px";
            exit.style.left="270px";
            exit.style.position="relative";
        }
      });

    });

    exit.addEventListener("click",()=>{
        value=0;
        diag.remove();
        div_operation.remove();
        div1.appendChild(div_background);
    });

    

});



//批量删除充电桩
const delBatch_url="/chargeports/delBatch";
let delBatch_data={
    ids:[]
}
const delBatch=document.querySelector("#delBatch");
delBatch.addEventListener("click",()=>{
    div_background.remove();
    div1.appendChild(div_operation);
    if (value !== 0) {
        return;
    }
    value=1;

    const diag=document.createElement("dialog");
    const input=document.createElement("input");
    const exit=document.createElement("button");
    const submit=document.createElement("button");
    exit.textContent="x";
    exit.id="exit";
    exit.style.left="100px";
    exit.style.bottom="50px";
    exit.className="btn btn-primary";
    submit.className="btn btn-primary";
    submit.textContent="确认";
    submit.id="submit";
    submit.style.top="100px";
    submit.style.right="130px";
    diag.textContent="请输入要删除的充电桩的id(用英文逗号隔开):";
    
    diag.appendChild(input);
    diag.appendChild(exit);
    diag.appendChild(submit);
    div_operation.appendChild(diag);
    diag.show();



    submit.addEventListener("click",()=>{
        let s=input.value.split(",");
        for(let i=0;i<s.length;i++){
            delBatch_data.ids.push(parseInt(s[i]))
        }
        console.log(delBatch_data.ids);
        const res=send_data(delBatch_url,delBatch_data);
        res.then(res=>res.json())
        .then(all_data=>{
            if(all_data.code===200){
                diag.textContent="删除成功";
                diag.appendChild(exit);
                exit.style.top="-19px";
                exit.style.left="270px";
                exit.style.position="relative";
            }
            else{
                diag.textContent="删除失败";
                diag.appendChild(exit);
                exit.style.top="-19px";
                exit.style.left="270px";
                exit.style.position="relative";
            }

        });

    });

    exit.addEventListener("click",()=>{
        value=0;
        diag.remove();
        div_operation.remove();
        div1.appendChild(div_background);
    });

});