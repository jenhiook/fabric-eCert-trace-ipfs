<template>
  <div class="tech-verify-container">
    <h3>技术实体端核验工具</h3>
    <el-card class="verify-card">
      <el-form label-width="100px">
        <el-form-item label="证照溯源码">
          <el-input v-model="certId" placeholder="请输入证照溯源码" clearable @change="loadData" />
        </el-form-item>
        <el-form-item label="选择历史记录（自动填充交易哈希）">
          <el-select v-model="selectedRecord" placeholder="加载历史后可选" clearable filterable @change="onSelectRecord">
            <el-option
              v-for="(record, idx) in historyList"
              :key="idx"
              :label="`${record.eventType} - ${record.time}`"
              :value="record"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="交易哈希（TxID）">
          <el-input v-model="txId" placeholder="自动填充" readonly />
        </el-form-item>
        <el-form-item label="IPFS CID（从证照自动获取）">
          <el-input v-model="cid" placeholder="自动填充" readonly />
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
import { getFruitInfo, getHistory } from '@/api/trace'

export default {
  data() {
    return {
      certId: '',
      certDetail: null,        // 证照详情
      historyList: [],
      selectedRecord: null,
      txId: '',
      cid: '',
      verifyResult: null,
      loading: false
    }
  },
  methods: {
    async loadData() {
      if (!this.certId) {
        this.certDetail = null
        this.historyList = []
        this.selectedRecord = null
        this.txId = ''
        this.cid = ''
        return
      }
      this.loading = true
      try {
        // 加载证照详情
        const fd = new FormData()
        fd.append('traceability_code', this.certId)
        const certRes = await getFruitInfo(fd)
        if (certRes.code === 200) {
          this.certDetail = JSON.parse(certRes.data)
          // 提取 CID（去除 ipfs:// 前缀）
          const img = this.certDetail.farmer_input?.fa_certImage || ''
          this.cid = img.replace(/^ipfs:\/\//, '')
        } else {
          this.certDetail = null
          this.cid = ''
        }
        // 加载历史记录
        const historyRes = await getHistory(this.certId)
        let list = []
        if (historyRes.code === 200 && historyRes.data) {
          if (typeof historyRes.data === 'string') {
            list = JSON.parse(historyRes.data)
          } else {
            list = historyRes.data
          }
        }
        this.historyList = list
        this.selectedRecord = null
        this.txId = ''
      } catch (e) {
        this.$message.error('加载数据失败')
      } finally {
        this.loading = false
      }
    },
    onSelectRecord(record) {
      if (record) {
        this.txId = record.txid || ''
      } else {
        this.txId = ''
      }
    },
    async verify() {
      if (!this.selectedRecord) {
        this.$message.warning('请先选择一条历史记录')
        return
      }
      if (!this.cid) {
        this.$message.warning('该证照无CID，无法核验')
        return
      }
      // 核验：交易哈希匹配，且 CID 与证照详情中的 CID 一致
      const txMatch = (this.txId === this.selectedRecord.txid)
      const cidMatch = (this.cid === (this.certDetail?.farmer_input?.fa_certImage || '').replace(/^ipfs:\/\//, ''))
      if (txMatch && cidMatch) {
        this.verifyResult = { valid: true, message: '核验通过，该证据与链上记录一致。' }
        this.$message.success('核验通过')
      } else {
        let msg = '核验失败：'
        if (!txMatch) msg += '交易哈希不匹配；'
        if (!cidMatch) msg += ' CID不匹配；'
        this.verifyResult = { valid: false, message: msg }
        this.$message.error('核验失败')
      }
    }
  }
}
</script>

<style scoped>
.tech-verify-container {
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
