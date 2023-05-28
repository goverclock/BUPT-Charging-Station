const btn1=document.querySelector("#signup");
const btn2=document.querySelector("#signin");
const user_name=document.querySelector("#username");
const pwd=document.querySelector("#password");
const body=document.querySelector("body");
const para=document.createElement("p");
body.appendChild(para);
let  server_addr="http://localhost:8080";

login_url="/login/user"
register_url="/register/user";
let usr={
    username:"",
    password:""}
btn1.addEventListener("click",()=>{
    usr.username=user_name.value;
    usr.password=pwd.value;
    const response=send_data(register_url,usr);
    response.then(response=>response.json())
    .then(data=>{
      console.log(data);
      if(data.code===200){
        para.textContent="注册成功请登录";
        
      }
      else{
        para.textContent="注册失败用户名已存在";
      }
      
    })
    
    
    
});

btn2.addEventListener("click",()=>{
    usr.username=user_name.value;
    usr.password=pwd.value;
    const response=send_data(login_url,usr);
    response.then(response=>response.json())
    .then(data=>{
      console.log(data);
      if(data.code===200){
        response.then(responseObject=>{
          console.log(responseObject.headers.get("Authorization"));
          console.log(responseObject.headers.get("Content-Type"));
          console.log(responseObject.headers);
          localStorage.setItem('tokens',responseObject.headers.get("Authorization"));
        });
        localStorage.setItem('user_id',data.data.id);
        localStorage.setItem('username',user_name.value);
        console.log(localStorage.getItem("username"));
        window.location.href="ui.html";
      }
      else{
        para.textContent="登录失败,用户名或密码错误";
      }
      
    })

});

//向服务器发送数据
function send_data(part_url,object){
   
    url=server_addr+part_url;
    const res=fetch(url ,{
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
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