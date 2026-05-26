<template>
  <div class="uplink-container">
    <h3>电子证照上链</h3>
    <p>当前用户：{{ user.name }} | 角色：{{ user.role }}</p>

    <div class="form-layout">
      <!-- 左侧：证照原图上传 -->
      <div class="upload-column">
        <div class="upload-card">
          <h4 class="upload-title">证照原图上传</h4>

          <div class="upload-box-wrapper" @click="triggerUpload">
            <img
              :src="previewImage"
              alt="证照预览"
              class="preview-img"
            >
            <input
              ref="fileInput"
              type="file"
              accept="image/*"
              hidden
              @change="handleFileSelect"
            >
          </div>

          <p class="upload-tip">点击上传证件照，支持身份证/毕业证/驾驶证等电子版</p>
        </div>
      </div>

      <!-- 右侧：表单 -->
      <div class="form-column">
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="130px"
        >
          <el-form-item label="证件类型" prop="type">
            <el-select v-model="form.type" placeholder="请选择证件类型">
              <el-option label="居民身份证" value="居民身份证" />
              <el-option label="毕业证书" value="毕业证书" />
              <el-option label="驾驶证" value="驾驶证" />
              <el-option label="职业资格证" value="职业资格证" />
            </el-select>
          </el-form-item>

          <el-form-item :label="label1" prop="field1">
            <el-input v-model="form.field1" :placeholder="placeholder1" />
          </el-form-item>

          <el-form-item :label="label2" prop="field2">
            <el-select
              v-if="form.type === '居民身份证'"
              v-model="form.field2"
              placeholder="性别"
            >
              <el-option label="男" value="男" />
              <el-option label="女" value="女" />
            </el-select>
            <el-input
              v-else
              v-model="form.field2"
              :placeholder="placeholder2"
            />
          </el-form-item>

          <el-form-item :label="label3" prop="field3">
            <el-input v-model="form.field3" :placeholder="placeholder3" />
          </el-form-item>

          <el-form-item label="持证人姓名" prop="name">
            <el-input v-model="form.name" placeholder="请输入真实姓名" />
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="submitForm">确认上链</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <div v-if="traceabilityCode" class="code-display">
      <el-alert title="上链成功！" type="info" :closable="false">
        <template #default>
          您的电子证照已成功上链，溯源码为：<strong>{{ traceabilityCode }}</strong>
        </template>
      </el-alert>
    </div>
  </div>
</template>

<script>
import { uplink } from '@/api/trace'

export default {
  data() {
    return {
      formRef: null,
      form: {
        type: '居民身份证',
        field1: '',
        field2: '',
        field3: '',
        name: ''
      },
      label1: '身份证号',
      label2: '性别',
      label3: '联系电话',
      placeholder1: '18位身份证号',
      placeholder2: '请选择性别',
      placeholder3: '11位手机号',

      // ✅ 修复：移除固定必填校验，控制台不再报错
      rules: {
        type: [{ required: true, message: '请选择证件类型', trigger: 'change' }],
        name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
        field1: [],
        field2: [],
        field3: []
      },

      traceabilityCode: null,
      previewImage: require('@/assets/camera.jpg'),
      selectedFile: null
    }
  },

  computed: {
    user() {
      const storeUser = this.$store.state.user || {}
      return {
        name: storeUser.username || storeUser.name || 'test',
        role: storeUser.role || '个人用户'
      }
    }
  },

  watch: {
    'form.type'(t) {
      if (t === '居民身份证') {
        this.label1 = '身份证号'
        this.label2 = '性别'
        this.label3 = '联系电话'
        this.placeholder1 = '请输入18位身份证号'
        this.placeholder2 = '请选择性别'
        this.placeholder3 = '请输入11位手机号'
      } else if (t === '毕业证书') {
        this.label1 = '证书编号'
        this.label2 = '毕业院校'
        this.label3 = '毕业时间'
        this.placeholder1 = '请输入证书编号'
        this.placeholder2 = '请输入学校名称'
        this.placeholder3 = '格式：2025-06-30'
      } else if (t === '驾驶证') {
        this.label1 = '驾驶证证号'
        this.label2 = '准驾车型'
        this.label3 = '有效期至'
        this.placeholder1 = '请输入驾驶证号'
        this.placeholder2 = '如 C1 B2 A1'
        this.placeholder3 = '格式：2035-12-31'
      } else if (t === '职业资格证') {
        this.label1 = '资格证编号'
        this.label2 = '工种/专业'
        this.label3 = '发证日期'
        this.placeholder1 = '请输入资格证号'
        this.placeholder2 = '如：软件工程师、护士'
        this.placeholder3 = '格式：2025-05-15'
      }
    }
  },

  methods: {
    triggerUpload() {
      this.$refs.fileInput.click()
    },
    handleFileSelect(e) {
      const file = e.target.files[0]
      if (!file) return
      if (!file.type.startsWith('image/')) {
        this.$message.warning('请上传图片文件')
        return
      }
      this.selectedFile = file
      const reader = new FileReader()
      reader.onload = (event) => {
        this.previewImage = event.target.result
      }
      reader.readAsDataURL(file)
    },

    async submitForm() {
      // 手动校验，保留必填逻辑
      if (!this.form.field1 || !this.form.field2 || !this.form.field3) {
        this.$message.error('请填写完整信息！')
        return
      }
      await this.$refs.formRef.validate()

      const type = this.form.type
      const f1 = this.form.field1.trim()
      const f3 = this.form.field3.trim()

      if (type === '居民身份证') {
        if (!/^\d{17}[\dXx]$/.test(f1)) {
          this.$message.error('身份证号必须是18位！')
          return
        }
        if (!/^1[3-9]\d{9}$/.test(f3)) {
          this.$message.error('手机号必须是11位有效号码！')
          return
        }
      }

      if (type === '毕业证书' || type === '驾驶证' || type === '职业资格证') {
        if (!/^\d{4}-\d{2}-\d{2}$/.test(f3)) {
          this.$message.error('日期格式必须是：YYYY-MM-DD')
          return
        }
      }

      const fd = new FormData()
      fd.append('arg1', this.form.type)
      fd.append('arg2', this.form.field1)
      fd.append('arg3', this.form.field2)
      fd.append('arg4', this.form.field3)
      fd.append('arg5', this.form.name)
      if (this.selectedFile) {
        fd.append('cert_image', this.selectedFile)
      }

      const res = await uplink(fd)
      if (res.code === 200) {
        this.traceabilityCode = res.traceability_code
        this.$message.success(`上链成功！溯源码：${res.traceability_code}`)
      } else {
        this.$message.error(res.message || '上链失败')
      }
    }
  }
}
</script>

<style scoped>
.uplink-container {
  width: 900px;
  margin: 30px auto;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  background: #fff;
}
h3 {
  font-size: 18px;
  margin-bottom: 10px;
  text-align: center;
}
p {
  font-size: 14px;
  color: #666;
  text-align: center;
}

.form-layout {
  display: flex;
  gap: 40px;
  align-items: flex-start;
  margin-top: 20px;
}
.upload-column { flex: 1; display: flex; justify-content: center; }
.upload-card {
  background: #fafbfc;
  border-radius: 8px;
  padding: 20px;
  text-align: center;
}
.upload-title { margin: 0 0 15px 0; font-size: 16px; color: #333; }

.upload-box-wrapper {
  width: 140px;
  aspect-ratio: 3 / 4;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  cursor: pointer;
  background: #fff;
  margin: 0 auto;
}
.preview-img {
  width: 90%;
  height: auto;
  object-fit: contain;
  display: block;
}
.upload-tip { margin: 12px 0 0 0; font-size: 12px; color: #999; }
.form-column { flex: 1.5; }

.code-display {
  margin-top: 20px;
}
</style>

