<template>
  <div class="event-container">
    <h3>合规事件存证（链上不可篡改）</h3>
    <el-form :inline="true">
      <el-form-item label="证照溯源码">
        <el-input v-model="traceCode" placeholder="请输入证照溯源码" clearable />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
        <el-button @click="showAddDialog">添加事件（上链）</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="eventList" border stripe>
      <el-table-column prop="event_type" label="事件类型" width="120" />
      <el-table-column prop="event_time" label="事件时间" width="160" />
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column label="附件CID" width="150">
        <template slot-scope="scope">
          <el-button
            v-if="scope.row.attachment_cid"
            type="text"
            @click="showAttachment(scope.row.attachment_cid)"
          >
            {{ scope.row.attachment_cid.substring(0, 20) }}...
          </el-button>
          <span v-else>无</span>
        </template>
      </el-table-column>
      <el-table-column prop="record_timestamp" label="存证时间" width="160" />
      <el-table-column prop="txid" label="交易哈希" width="200" show-overflow-tooltip />
    </el-table>

    <!-- 添加事件对话框 -->
    <el-dialog title="添加合规事件" :visible.sync="dialogVisible" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="证照溯源码" required>
          <el-input v-model="form.traceability_code" disabled />
        </el-form-item>
        <el-form-item label="事件类型" required>
          <el-select v-model="form.event_type" placeholder="请选择">
            <el-option label="年检通过" value="年检通过" />
            <el-option label="资质更新" value="资质更新" />
            <el-option label="异常报告" value="异常报告" />
            <el-option label="整改完成" value="整改完成" />
          </el-select>
        </el-form-item>
        <el-form-item label="事件时间" required>
          <el-date-picker v-model="form.event_time" type="datetime" placeholder="选择时间" value-format="yyyy-MM-dd HH:mm:ss" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input type="textarea" v-model="form.description" />
        </el-form-item>
        <el-form-item label="附件CID">
          <el-input
            v-model="form.attachment_cid"
            placeholder="请输入IPFS CID（可通过上传附件获取或直接粘贴）"
          />
          <el-button @click="uploadToIPFS" size="small">上传附件</el-button>
          <div style="font-size:12px; color:#909399; margin-top:4px;">
            支持手动输入CID，或点击上传按钮自动填充。
          </div>
        </el-form-item>
        <el-form-item label="指纹">
          <el-input v-model="form.attachment_fingerprint" placeholder="可选，用于校验" />
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAdd">确认上链</el-button>
      </span>
    </el-dialog>

    <!-- 附件预览弹窗：只提供链接 -->
    <el-dialog title="附件预览" :visible.sync="attachmentDialogVisible" width="50%">
      <div v-if="attachmentUrl">
        <p>请点击下方链接查看附件内容：</p>
        <el-link :href="attachmentUrl" target="_blank" type="primary" icon="el-icon-link">{{ attachmentUrl }}</el-link>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { addComplianceEvent, getComplianceEvents } from '@/api/enterprise'
import request from '@/utils/request'

export default {
  data() {
    return {
      traceCode: '',
      eventList: [],
      dialogVisible: false,
      form: {
        traceability_code: '',
        event_type: '',
        event_time: '',
        description: '',
        attachment_cid: '',
        attachment_fingerprint: ''
      },
      attachmentDialogVisible: false,
      attachmentUrl: ''
    }
  },
  methods: {
    async loadList() {
      if (!this.traceCode) {
        this.$message.warning('请输入证照溯源码')
        return
      }
      const res = await getComplianceEvents(this.traceCode)
      if (res.code === 200) {
        // 按事件时间降序（最新的事件在前）
        this.eventList = res.data.sort((a, b) => (b.event_time || '').localeCompare(a.event_time || ''))
      }
    },
    showAddDialog() {
      if (!this.traceCode) {
        this.$message.warning('请先输入证照溯源码')
        return
      }
      this.form.traceability_code = this.traceCode
      this.dialogVisible = true
      // 清空上次的CID和指纹
      this.form.attachment_cid = ''
      this.form.attachment_fingerprint = ''
    },
    async uploadToIPFS() {
      const input = document.createElement('input')
      input.type = 'file'
      input.onchange = async (e) => {
        const file = e.target.files[0]
        const fd = new FormData()
        fd.append('file', file)
        try {
          const res = await request.post('/upload/ipfs', fd)
          if (res.code === 200) {
            this.form.attachment_cid = res.cid
            this.$message.success('上传成功，CID已填入')
          } else {
            this.$message.error('上传失败: ' + (res.message || '未知错误'))
          }
        } catch (err) {
          this.$message.error('上传请求失败，请检查后端服务')
        }
      }
      input.click()
    },
    isValidCID(cid) {
      if (!cid) return true // 允许为空
      // 简单校验 IPFS CID v0 (Qm开头，46字符) 或 v1 (b开头，59字符)
      const cidV0Regex = /^Qm[1-9A-HJ-NP-Za-km-z]{44}$/
      const cidV1Regex = /^b[A-Za-z2-7]{58}$/
      return cidV0Regex.test(cid) || cidV1Regex.test(cid)
    },
    async submitAdd() {
      if (!this.form.event_type || !this.form.event_time) {
        this.$message.warning('请填写完整信息')
        return
      }
      if (this.form.attachment_cid && !this.isValidCID(this.form.attachment_cid)) {
        this.$message.error('附件CID格式不正确，应为IPFS内容标识符（如 Qm...），请检查或使用上传按钮获取正确的CID。')
        return
      }
      const res = await addComplianceEvent(this.form)
      if (res.code === 200) {
        this.$message.success('上链成功，交易哈希：' + (res.txid || '无'))
        this.dialogVisible = false
        this.loadList()
      } else {
        this.$message.error(res.message || '添加失败')
      }
    },
    showAttachment(cid) {
      if (!cid) return
      this.attachmentUrl = `http://127.0.0.1:8080/ipfs/${cid}`
      this.attachmentDialogVisible = true
    }
  }
}
</script>

<style scoped>
.event-container {
  padding: 20px;
}
</style>
