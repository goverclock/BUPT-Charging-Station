const btn1=document.querySelector("#signup");
const btn2=document.querySelector("#signin");
const user_name=document.querySelector("#username");
const pwd=document.querySelector("#password");
const body=document.querySelector("body");
const para=document.createElement("p");
body.appendChild(para);

login_url="/user/login"
let usr={
    username:"",
    password:""}
btn1.addEventListener("click",()=>{
    usr.username=user_name.value;
    usr.password=pwd.value;
    const response=send_data(login_url,usr);
    response.then(response=>response.json())
    .then(data=>{
      console.log(data);
      if(data.code===200){
        window.location.href="login.html";
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
        window.location.href="ui.html";
      }
      else{
        //window.location.href="login.html";
        para.textContent="登录失败,用户名或密码错误";
      }
      
    })

});

//向服务器发送数据
function send_data(part_url,object){
    server_addr="http://localhost:8080";
    url=server_addr+part_url;
    const res=fetch(url , {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(object)
      });
      return res;
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