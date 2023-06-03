const btn1=document.querySelector("#signup");
const btn2=document.querySelector("#signin");
const user_name=document.querySelector("#username");
const pwd=document.querySelector("#password");
const body=document.querySelector("body");
const para=document.createElement("p");
const ip=document.querySelector("#ip_address")
const port=document.querySelector("#port");
body.appendChild(para);
let  server_addr="http://localhost:8080";

login_url="/login/user"
register_url="/register/user";
let usr={
    username:"",
    password:""}
btn1.addEventListener("click",()=>{
    server_addr="http://"+ip.value+":"+port.value;
    console.log(server_addr);
    usr.username=user_name.value;
    usr.password=pwd.value;
    const response=send_data(register_url,usr);
    response.then(response=>response.json())
    .then(data=>{
      console.log(data);
      if(data.code===200){
        para.textContent=data.msg;
        
      }
      else{
        para.textContent=data.msg;
      }
      
    })
});

btn2.addEventListener("click",()=>{
    usr.username=user_name.value;
    usr.password=pwd.value;
    server_addr="http://"+ip.value+":"+port.value;
    console.log(server_addr);
    const response=send_data(login_url,usr);
    response.then(response=>response.json())
    .then(data=>{
      console.log(data);
      if(data.code===200){
        response.then(responseObject=>{
          console.log(responseObject.headers.get("Authorization"));
          console.log(responseObject.headers.get("Content-Type"));
          console.log(responseObject.headers);
          localStorage.setItem('tokens',data.data.token);
          localStorage.setItem("address",server_addr);
        });
        if(data.data.user_type===0){
           localStorage.setItem('user_id',data.data.user_id);
           localStorage.setItem('username',user_name.value);
           console.log(localStorage.getItem("username"));
           console.log(localStorage.getItem("user_id"));
           window.location.href="ui.html";
        }
        else{
           localStorage.setItem('admin_id',data.data.user_id);
           localStorage.setItem('adminname',user_name.value);
           console.log(localStorage.getItem("adminname"));
           window.location.href="admin.html";

        }
      }
      else{
        para.textContent=data.msg;
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