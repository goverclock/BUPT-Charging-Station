
    const input1 = document.getElementById("exampleInputEmail1");
    const input2 = document.getElementById("password");
    const btn1=document.getElementById("sign in");
    const btn2=document.getElementById("sign up");
    const para = document.getElementById("output");
    const p=document.getElementById("1");
    const div1=document.getElementById("div1");
    const div2=document.getElementById("div2");
    const div3=document.getElementById("div3");
    
    


    btn1.addEventListener("click", () => {
        const yonghu=input1.value;
        const pwd=input2.value;
        if(yonghu !=="" && pwd!==""){
            para.textContent=yonghu+"注册成功,请登录";
            input1.value=" ";
            input2.value="";
        }
        else{
            if(yonghu===""){
                para.textContent="用户名不能为空";
            }
            else{
                para.textContent="密码不能为空";

            }
        }
  
});
    btn2.addEventListener("click", () => {
        const yonghu=input1.value;
        const pwd=input2.value;
        //if(yonghu !=="" && pwd!==""){
            input1.remove();
            input2.remove();
            btn1.remove();
            btn2.remove();
            para.remove();
            p.remove();
            div1.remove();
            div3.remove();
            div2.remove();
            init();
        //}

});
function init(){

  let newP = document.createElement("h1");
  let btnlist=document.createElement("div");
  let start_charge=document.createElement("button");
  let mod_req=document.createElement("button");
  let end_charge=document.createElement("button");
  let exit=document.createElement("button");

  newP.textContent="欢迎使用自助充电系统";
  btnlist.textContent="请选择你的操作:";
  start_charge.className="btn btn-primary";
  mod_req.className="btn btn-primary";
  end_charge.className="btn btn-primary";
  exit.className="btn btn-primary";
  btnlist.className="btnlist";
  
  start_charge.textContent="开始充电";
  mod_req.textContent="调整方式";
  end_charge.textContent="停止充电";
  exit.textContent="退出程序";

  document.body.appendChild(newP);
  document.body.appendChild(btnlist);
  document.body.appendChild(start_charge);
  document.body.appendChild(mod_req);
  document.body.appendChild(end_charge);
  document.body.appendChild(exit);

  

  exit.addEventListener("click",exit_);

}
function start_charge_(){

}
function end_charge_(){

}

function mod_req_(){


}
function exit_(){
  window.close();

}
