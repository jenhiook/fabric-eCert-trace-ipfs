<template>
  <div class="enterprise-verify-container">
    <h3>企业端证据核验工具</h3>
    <el-card class="verify-card">
      <el-form label-width="100px">
        <el-form-item label="证照溯源码">
          <el-input v-model="certId" placeholder="请输入证照溯源码" clearable />
        </el-form-item>
        <el-form-item label="交易哈希（TxID）">
          <el-input v-model="txId" placeholder="请从证照详情页完整复制交易哈希（64位）" clearable />
        </el-form-item>
        <el-form-item label="IPFS CID">
          <el-input v-model="cid" placeholder="请从证照详情页复制IPFS CID" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="verify">开始核验</el-button>
        </el-form-item>
      </el-form>
      <div v-if="verifyResult" class="verify-result">
        <el-alert :title="verifyResult.message" :type="verifyResult.valid ? 'success' : 'error'" :closable="false" />
      </div>
    </el-card>
  </div>
</template>

<script>
import { getHistory } from '@/api/trace'

export default {
  data() {
    return {
      certId: '',
      txId: '',
      cid: '',
      verifyResult: null,
      loading: false
    }
  },
  methods: {
    async verify() {
      if (!this.certId || !this.txId || !this.cid) {
        this.$message.warning('请填写完整信息')
        return
      }
      this.loading = true
      this.verifyResult = null
      try {
        const res = await getHistory(this.certId)
        let history = []
        if (res.code === 200 && res.data) {
          if (typeof res.data === 'string') {
            history = JSON.parse(res.data)
          } else {
            history = res.data
          }
        }
        const matched = history.some(record => record.txid === this.txId && record.cid === this.cid)
        if (matched) {
          this.verifyResult = { valid: true, message: '核验通过，该证据与链上记录一致。' }
          this.$message.success('核验通过')
        } else {
          this.verifyResult = { valid: false, message: '核验失败：交易哈希或CID与链上记录不匹配。请检查是否完整复制，或该证照历史记录不完整（建议用新证照测试）。' }
          this.$message.error('核验失败')
        }
      } catch (error) {
        console.error(error)
        this.$message.error('获取证照历史失败，请检查溯源码或网络')
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.enterprise-verify-container {
  padding: 20px;
}
.verify-card {
  max-width: 700px;
  margin: 20px auto;
}
.verify-result {
  margin-top: 20px;
}
</style>
