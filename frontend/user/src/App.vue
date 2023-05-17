<script>
export default {
  data() {
    return {
      server_addr: "http://localhost:8080"    // private IP address is also ok, e.g. http://10.128.135.80:8080
    }
  },
  methods: {
    login() {
      const username = document.querySelector('input[name="user_name"]').value
      const password = document.querySelector('input[name="user_password"]').value
      const log_info = { username, password }
      console.log(log_info)
      fetch(this.server_addr + '/user/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(log_info)
      })
        .then(response => response.json())
        .then(log_info => {
          console.log(log_info);
        })
        .catch(error => {
          console.error(error);
        });
    }
  }
}
</script>

<template>
  <div>
    <h3>
      欢迎
      <small class="text-muted">登录充电系统</small>
    </h3>

    <div>
      <label class="form-label mt-4" for="exampleInputEmail1">用户名</label>
      <input type="email" name="user_name" class="form-control" id="exampleInputEmail1" aria-describedby="emailHelp"
        placeholder="输入邮箱地址">
    </div>

    <div>
      <label for="password" class="form-label mt-4">密码</label>
      <input type="password" name="user_password" class="form-control" placeholder="输入密码" id="password">
      <br>
    </div>

    <button id="sign in" name="sign_in" class="btn btn-primary" @click="login">
      登录
    </button>
    <button id="sign up" name="sign_up" class="btn btn-secondary">
      注册
    </button>

  </div>
</template>


<style src="./assets/style.css"></style>
