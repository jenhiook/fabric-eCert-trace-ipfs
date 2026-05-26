<template>
  <div class="govt-audit-container">
    <h3>政务部门 — 证照审核</h3>
    <p class="role-info">当前角色：政务部门 | 待审核：{{ pendingList.length }} 条</p>

    <!-- 统计卡片区 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="12">
        <el-card
          class="stat-card"
          :class="{ active: activeTab === 'pending' }"
          shadow="hover"
          @click.native="activeTab = 'pending'"
        >
          <div class="stat-icon"><i class="el-icon-warning-outline" /></div>
          <div class="stat-info">
            <div class="stat-number">{{ pendingList.length }}</div>
            <div class="stat-title">未审核证照</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card
          class="stat-card"
          :class="{ active: activeTab === 'audited' }"
          shadow="hover"
          @click.native="activeTab = 'audited'"
        >
          <div class="stat-icon"><i class="el-icon-circle-check" /></div>
          <div class="stat-info">
            <div class="stat-number">{{ auditedList.length }}</div>
            <div class="stat-title">已审核证照</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 待审核证照列表 -->
    <el-card v-if="activeTab === 'pending'" class="list-card">
      <div slot="header">
        <span>待审核证照（{{ pendingList.length }} 条）</span>
      </div>
      <el-table :data="pagedPendingList" border stripe>
        <el-table-column prop="traceability_code" label="溯源码" width="180" />
        <el-table-column prop="farmer_input.fa_fruitName" label="证照类型" />
        <el-table-column prop="farmer_input.fa_farmerName" label="持证人" />
        <el-table-column prop="farmer_input.fa_origin" label="证照编号" />
        <el-table-column prop="farmer_input.fa_timestamp" label="上链时间" width="160" />
        <el-table-column label="操作" width="100">
          <template slot-scope="scope">
            <el-button type="primary" size="small" @click="showAuditDialog(scope.row)">审核</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="pendingList.length"
        :page-size="pendingPageSize"
        :current-page.sync="pendingCurrentPage"
        style="margin-top:15px; text-align:right;"
        @current-change="handlePendingPageChange"
      />
    </el-card>

    <!-- 已审核证照列表 -->
    <el-card v-if="activeTab === 'audited'" class="list-card">
      <div slot="header">
        <span>已审核证照（{{ auditedList.length }} 条）</span>
      </div>
      <el-table :data="pagedAuditedList" border stripe>
        <el-table-column prop="traceability_code" label="溯源码" width="180" />
        <el-table-column prop="farmer_input.fa_fruitName" label="证照类型" />
        <el-table-column prop="farmer_input.fa_farmerName" label="持证人" />
        <el-table-column prop="farmer_input.fa_origin" label="证照编号" />
        <el-table-column prop="factory_input.fac_timestamp" label="审核时间" width="160" />
        <el-table-column label="审核结果" width="100">
          <template slot-scope="scope">
            <el-tag :type="scope.row.factory_input.fac_contactNumber === '通过' ? 'success' : 'danger'">
              {{ scope.row.factory_input.fac_contactNumber === '通过' ? '通过' : '驳回' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template slot-scope="scope">
            <el-button size="small" @click="viewDetail(scope.row)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="auditedList.length"
        :page-size="auditedPageSize"
        :current-page.sync="auditedCurrentPage"
        style="margin-top:15px; text-align:right;"
        @current-change="handleAuditedPageChange"
      />
    </el-card>

    <!-- 审核对话框 -->
    <el-dialog title="证照审核" :visible.sync="dialogVisible" width="500px">
      <el-form :model="auditForm" label-width="100px">
        <el-form-item label="证照类型">
          <span>{{ currentCert && currentCert.farmer_input && currentCert.farmer_input.fa_fruitName }}</span>
        </el-form-item>
        <el-form-item label="持证人">
          <span>{{ currentCert && currentCert.farmer_input && currentCert.farmer_input.fa_farmerName }}</span>
        </el-form-item>
        <el-form-item label="证照编号">
          <span>{{ currentCert && currentCert.farmer_input && currentCert.farmer_input.fa_origin }}</span>
        </el-form-item>
        <el-form-item label="部门名称" required>
          <el-input v-model="auditForm.deptName" placeholder="请输入部门名称" />
        </el-form-item>
        <el-form-item label="部门代码" required>
          <el-input v-model="auditForm.deptCode" placeholder="请输入部门代码" />
        </el-form-item>
        <el-form-item label="审核人" required>
          <el-input v-model="auditForm.auditorName" placeholder="请输入审核人姓名" />
        </el-form-item>
        <el-form-item label="审核结果" required>
          <el-radio-group v-model="auditForm.auditResult">
            <el-radio label="通过">通过</el-radio>
            <el-radio label="驳回">驳回</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAudit">提交</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { getAllFruitInfo, govtAudit } from '@/api/trace'

export default {
  data() {
    return {
      allCerts: [],
      activeTab: 'pending',
      dialogVisible: false,
      currentCert: null,
      auditForm: {
        deptName: '',
        deptCode: '',
        auditorName: '',
        auditResult: '通过'
      },
      pendingCurrentPage: 1,
      pendingPageSize: 10,
      auditedCurrentPage: 1,
      auditedPageSize: 10
    }
  },
  computed: {
    pendingList() {
      return this.allCerts.filter(cert => {
        const farmer = cert.farmer_input || {}
        const govt = cert.factory_input || {}
        return farmer.fa_txid && (!govt.fac_txid || govt.fac_txid === '')
      })
    },
    auditedList() {
      return this.allCerts.filter(cert => {
        const govt = cert.factory_input || {}
        return govt.fac_txid && govt.fac_txid !== ''
      })
    },
    pagedPendingList() {
      const start = (this.pendingCurrentPage - 1) * this.pendingPageSize
      return this.pendingList.slice(start, start + this.pendingPageSize)
    },
    pagedAuditedList() {
      const start = (this.auditedCurrentPage - 1) * this.auditedPageSize
      return this.auditedList.slice(start, start + this.auditedPageSize)
    }
  },
  mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      const res = await getAllFruitInfo()
      if (res.code === 200 && res.data) {
        this.allCerts = JSON.parse(res.data)
      }
    },
    showAuditDialog(row) {
      this.currentCert = row
      this.dialogVisible = true
      this.auditForm = {
        deptName: '',
        deptCode: '',
        auditorName: '',
        auditResult: '通过'
      }
    },
    async submitAudit() {
      console.log('提交审核, 表单数据:', this.auditForm)
      if (!this.auditForm.deptName || !this.auditForm.deptCode || !this.auditForm.auditorName) {
        this.$message.warning('请填写完整信息')
        return
      }
      const fd = new FormData()
      fd.append('traceability_code', this.currentCert.traceability_code)
      fd.append('deptName', this.auditForm.deptName)
      fd.append('deptCode', this.auditForm.deptCode)
      fd.append('auditorName', this.auditForm.auditorName)
      fd.append('auditResult', this.auditForm.auditResult)
      try {
        const res = await govtAudit(fd)
        if (res.code === 200) {
          this.$message.success('审核成功')
          this.dialogVisible = false
          this.loadData()
        } else {
          this.$message.error(res.message || '审核失败')
        }
      } catch (error) {
        console.error(error)
        this.$message.error('审核请求异常')
      }
    },
    viewDetail(row) {
      this.$router.push({ path: '/trace', query: { code: row.traceability_code }})
    },
    handlePendingPageChange(page) {
      this.pendingCurrentPage = page
    },
    handleAuditedPageChange(page) {
      this.auditedCurrentPage = page
    }
  }
}
</script>

<style lang="scss" scoped>
.govt-audit-container {
  padding: 20px;
  background: #f0f2f5;
  min-height: 100vh;
}
.role-info {
  margin-bottom: 20px;
  font-size: 14px;
  color: #606266;
}
.stat-cards {
  margin-bottom: 20px;
}
.stat-card {
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
  border-radius: 16px;
  &.active {
    border-color: #409eff;
    box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
  }
  .stat-icon {
    font-size: 48px;
    color: #409eff;
    margin-bottom: 10px;
  }
  .stat-number {
    font-size: 32px;
    font-weight: bold;
    color: #303133;
  }
  .stat-title {
    font-size: 14px;
    color: #909399;
    margin-top: 8px;
  }
}
.list-card {
  border-radius: 16px;
  overflow: hidden;
}
</style>
