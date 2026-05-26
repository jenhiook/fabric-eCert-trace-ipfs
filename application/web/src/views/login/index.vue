<template>
  <div class="login-container">
    <div class="login-card">
      <!-- 标题 -->
      <div class="title-area">
        <h2 class="sys-title">基于区块链的</h2>
        <h2 class="sys-title">电子证照存储与溯源系统</h2>
        <p class="sub-title">组长：23级信息安全 吴晓汀 &nbsp;|&nbsp; 组员：翁茹芳</p>
      </div>

      <el-divider />

      <!-- 登录 -->
      <div v-show="isLoginPage">
        <el-form ref="loginForm" :model="loginForm" :rules="loginRules" size="medium">
          <el-form-item prop="username">
            <el-input v-model="loginForm.username" placeholder="请输入账号" prefix-icon="el-icon-user" />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              :type="passwordType"
              placeholder="请输入密码"
              prefix-icon="el-icon-lock"
              @keyup.enter.native="handleLogin"
            >
              <i
                slot="suffix"
                :class="passwordType === 'password' ? 'el-icon-view' : 'el-icon-loading'"
                style="cursor:pointer;line-height:40px"
                @click="showPwd"
              />
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="medium" :loading="loading" style="width:100%" @click="handleLogin">登 录</el-button>
          </el-form-item>
          <el-form-item style="text-align:center;margin-bottom:0">
            <el-button type="text" @click="handleRegister">没有账号？去注册</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 注册 -->
      <div v-show="!isLoginPage">
        <el-form :model="registerForm" size="medium">
          <el-form-item>
            <el-input v-model="registerForm.username" placeholder="请输入账号" prefix-icon="el-icon-user" />
          </el-form-item>
          <el-form-item>
            <el-input v-model="registerForm.password" :type="passwordType" placeholder="请输入密码" prefix-icon="el-icon-lock" />
          </el-form-item>
          <el-form-item>
            <el-input v-model="registerForm.password2" :type="passwordType" placeholder="请再次输入密码" prefix-icon="el-icon-lock" />
          </el-form-item>
          <el-form-item>
            <el-select v-model="registerForm.userType" placeholder="请选择角色" style="width:100%">
              <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="medium" :loading="loading" style="width:100%" @click="submitRegister">提交注册</el-button>
          </el-form-item>
          <el-form-item style="text-align:center;margin-bottom:0">
            <el-button type="text" @click="handleRegister">已有账号？返回登录</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Login',
  data() {
    return {
      loginForm: { username: '', password: '' },
      loginRules: {
        username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
      },
      loading: false,
      passwordType: 'password',
      redirect: undefined,
      isLoginPage: true,
      registerForm: { username: '', password: '', password2: '', userType: '' },
      options: [
        { value: '个人用户', label: '个人用户' },
        { value: '政务部门', label: '政务部门' },
        { value: '企业组织', label: '企业组织' },
        { value: '技术支撑实体', label: '技术支撑实体' }
      ]
    }
  },
  watch: {
    $route: {
      handler: route => { this.redirect = (route.query || {}).redirect },
      immediate: true
    }
  },
  methods: {
    showPwd() {
      this.passwordType = this.passwordType === 'password' ? 'text' : 'password'
    },

    handleLogin() {
      this.$refs.loginForm.validate(valid => {
        if (!valid) return
        this.loading = true
        this.$store.dispatch('user/login', this.loginForm).then(() => {
          const role = this.$store.state.user.userType || ''
          const roleMap = {
            '个人用户': '/uplink',
            '政务部门': '/govt-audit',
            '企业组织': '/enterprise',
            '技术支撑实体': '/tech-support'
          }
          this.$router.push({ path: this.redirect || roleMap[role] || '/home' })
        }).finally(() => { this.loading = false })
      })
    },

    handleRegister() {
      this.isLoginPage = !this.isLoginPage
    },
    submitRegister() {
      if (this.registerForm.password !== this.registerForm.password2) {
        return this.$message.error('两次密码不一致')
      }
      if (!this.registerForm.userType) {
        return this.$message.error('请选择角色')
      }
      const loading = this.$loading({ lock: true, text: '注册中...', background: 'rgba(0,0,0,0.7)' })
      this.$store.dispatch('user/register', this.registerForm).then(res => {
        this.$message.success('注册成功！交易ID：' + res.txid)
        this.handleRegister()
      }).catch(err => {
        this.$message.error(err.message || '注册失败')
      }).finally(() => loading.close())
    }
  }
}
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  width: 100%;
  background-image: url("../../assets/IMG_20260515_111727.jpg");
  background-size: cover;
  background-position: center;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-card {
  width: 560px;
  padding: 50px 50px 36px;
  /* 登录卡片本身也略微调整为更通透的玻璃质感 */
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 0 10px 50px rgba(0, 0, 0, 0.3);
}

.title-area {
  text-align: center;
  margin-bottom: 14px;

  .sys-title {
    font-size: 26px;
    font-weight: 700;
    color: #1a1a2e;
    margin: 0 0 2px;
    letter-spacing: 2px;
  }

  .sub-title {
    font-size: 14px;
    color: #999;
    margin-top: 12px;
  }
}

/* ================= 核心修复：按钮样式 ================= */

/* 主按钮：天青色 + 半透明玻璃质感 */
::v-deep .login-card .el-button--primary {
  background-color: rgba(90, 158, 155, 0.5) !important; /* 半透明天青色 */
  border-color: rgba(90, 158, 155, 0.5) !important;
  backdrop-filter: blur(4px); /* 玻璃毛玻璃效果 */
  box-shadow: 0 4px 12px rgba(90, 158, 155, 0.3); /* 发光阴影 */
  transition: all 0.3s ease;
}

::v-deep .login-card .el-button--primary:hover {
  background-color: rgb(129, 179, 202) !important; /* 悬浮时变实色 */
  border-color: rgb(144, 173, 218) !important;
  box-shadow: 0 6px 16px rgba(90, 158, 155, 0.4);
  transform: translateY(-1px);
}

::v-deep .login-card .el-button--primary:active {
  background-color: rgba(143, 203, 229, 0.9) !important;
  border-color: rgba(143, 198, 227, 0.9) !important;
  transform: translateY(1px);
}

/* 文字链接按钮：颜色跟随背景调性 */
::v-deep .login-card .el-button--text {
  color: #87bce7 !important; /* 改为和主色调一致的色 */
  font-weight: 500;
  transition: all 0.3s ease;
}

::v-deep .login-card .el-button--text:hover {
  color: #91b2dc !important;
  text-shadow: 0 0 8px rgba(90, 158, 155, 0.3);
}

/* ==================================================== */

/* 输入框加大 */
::v-deep .el-input__inner {
  height: 50px;
  font-size: 16px;
  border-radius: 8px;
  /* 输入框也可以稍微加一点透明感 */
  background-color: rgba(255, 255, 255, 0.7);
}
::v-deep .el-input__inner:focus {
  border-color: #7da8d1;
}

/* 下拉框同样加大 */
::v-deep .el-select .el-input__inner {
  height: 50px;
  font-size: 16px;
}

/* 按钮加大 */
::v-deep .el-button--medium {
  height: 50px;
  font-size: 17px;
  letter-spacing: 8px;
  border-radius: 8px;
}

::v-deep .el-divider {
  margin: 16px 0 22px;
}
</style>
