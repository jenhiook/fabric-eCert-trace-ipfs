<template>
  <div class="enterprise-container">
    <h3>企业组织 — 证照使用备案</h3>
    <p class="role-info">当前角色：企业组织 | 待备案：{{ pendingList.length }} 条</p>

    <el-row :gutter="20" class="stat-cards">
      <el-col :span="12">
        <el-card class="stat-card" :class="{ active: activeTab === 'pending' }" shadow="hover" @click.native="activeTab = 'pending'">
          <div class="stat-icon"><i class="el-icon-warning-outline" /></div>
          <div class="stat-info">
            <div class="stat-number">{{ pendingList.length }}</div>
            <div class="stat-title">待备案证照</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="stat-card" :class="{ active: activeTab === 'used' }" shadow="hover" @click.native="activeTab = 'used'">
          <div class="stat-icon"><i class="el-icon-circle-check" /></div>
          <div class="stat-info">
            <div class="stat-number">{{ usedList.length }}</div>
            <div class="stat-title">已备案证照</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 待备案证照列表 -->
    <el-card v-if="activeTab === 'pending'" class="list-card">
      <div slot="header"><span>待备案证照（{{ pendingList.length }} 条）</span></div>
      <el-table v-loading="loading" :data="pagedPendingList" border stripe fit>
        <el-table-column prop="traceability_code" label="溯源码" min-width="180" />
        <el-table-column label="证照类型" min-width="110">
          <template slot-scope="scope">{{ scope.row.farmer_input.fa_fruitName }}</template>
        </el-table-column>
        <el-table-column label="持证人" min-width="100">
          <template slot-scope="scope">{{ scope.row.farmer_input.fa_farmerName }}</template>
        </el-table-column>
        <el-table-column label="证照编号" min-width="200">
          <template slot-scope="scope">{{ scope.row.farmer_input.fa_origin }}</template>
        </el-table-column>
        <el-table-column label="审核状态" min-width="100">
          <template slot-scope="scope"><el-tag type="success">已审核</el-tag></template>
        </el-table-column>
        <el-table-column label="操作" min-width="100">
          <template slot-scope="scope">
            <el-button type="primary" size="small" @click="openUse(scope.row)">备案使用</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination background layout="total, prev, pager, next" :total="pendingList.length" :page-size="pendingPageSize" :current-page.sync="pendingCurrentPage" style="margin-top:15px; text-align:right;" @current-change="handlePendingPageChange" />
    </el-card>

    <!-- 已备案证照列表 -->
    <el-card v-if="activeTab === 'used'" class="list-card">
      <div slot="header"><span>已备案证照（{{ usedList.length }} 条）</span></div>
      <el-table :data="pagedUsedList" border stripe fit>
        <el-table-column prop="traceability_code" label="溯源码" min-width="180" />
        <el-table-column label="证照类型" min-width="110">
          <template slot-scope="scope">{{ scope.row.farmer_input.fa_fruitName }}</template>
        </el-table-column>
        <el-table-column label="持证人" min-width="100">
          <template slot-scope="scope">{{ scope.row.farmer_input.fa_farmerName }}</template>
        </el-table-column>
        <el-table-column label="企业名称" min-width="140">
          <template slot-scope="scope">{{ scope.row.driver_input.dr_name }}</template>
        </el-table-column>
        <el-table-column label="信用代码" min-width="180">
          <template slot-scope="scope">{{ scope.row.driver_input.dr_age }}</template>
        </el-table-column>
        <el-table-column label="使用目的" min-width="140">
          <template slot-scope="scope">{{ scope.row.driver_input.dr_phone }}</template>
        </el-table-column>
        <el-table-column label="经办人" min-width="100">
          <template slot-scope="scope">{{ scope.row.driver_input.dr_transport }}</template>
        </el-table-column>
        <el-table-column label="备案时间" min-width="180">
          <template slot-scope="scope">{{ scope.row.driver_input.dr_timestamp }}</template>
        </el-table-column>
        <el-table-column label="操作" min-width="100">
          <template slot-scope="scope">
            <el-button size="small" @click="viewDetail(scope.row)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination background layout="total, prev, pager, next" :total="usedList.length" :page-size="usedPageSize" :current-page.sync="usedCurrentPage" style="margin-top:15px; text-align:right;" @current-change="handleUsedPageChange" />
    </el-card>

    <!-- 备案对话框 -->
    <el-dialog title="证照使用备案" :visible.sync="dialogVisible" width="500px">
      <el-descriptions v-if="currentCert" :column="1" border>
        <el-descriptions-item label="溯源码">{{ currentCert.traceability_code }}</el-descriptions-item>
        <el-descriptions-item label="证照类型">{{ (currentCert.farmer_input || {}).fa_fruitName }}</el-descriptions-item>
        <el-descriptions-item label="持证人">{{ (currentCert.farmer_input || {}).fa_farmerName }}</el-descriptions-item>
      </el-descriptions>
      <el-form ref="useForm" :model="use" :rules="rules" label-width="100px" style="margin-top:20px">
        <el-form-item label="企业名称" prop="companyName">
          <el-select v-model="use.companyName" placeholder="请选择" style="width:100%">
            <el-option label="阿里巴巴集团" value="阿里巴巴集团" />
            <el-option label="腾讯科技" value="腾讯科技" />
            <el-option label="华为技术" value="华为技术" />
            <el-option label="字节跳动" value="字节跳动" />
            <el-option label="百度在线" value="百度在线" />
            <el-option label="京东集团" value="京东集团" />
          </el-select>
        </el-form-item>
        <el-form-item label="信用代码" prop="companyCode">
          <el-input v-model="use.companyCode" placeholder="统一社会信用代码" />
        </el-form-item>
        <el-form-item label="使用目的" prop="usePurpose">
          <el-select v-model="use.usePurpose" placeholder="请选择" style="width:100%">
            <el-option label="入职资格审查" value="入职资格审查" />
            <el-option label="项目招投标" value="项目招投标" />
            <el-option label="资质年检备案" value="资质年检备案" />
            <el-option label="合同签署认证" value="合同签署认证" />
            <el-option label="安全检查审核" value="安全检查审核" />
          </el-select>
        </el-form-item>
        <el-form-item label="经办人" prop="operator">
          <el-input v-model="use.operator" placeholder="经办人姓名" />
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="doUse">确认备案</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { getAllFruitInfo } from '@/api/trace'
import request from '@/utils/request'

export default {
  data() {
    return {
      allCerts: [],
      loading: false,
      activeTab: 'pending',
      dialogVisible: false,
      currentCert: null,
      use: { companyName: '', companyCode: '', usePurpose: '', operator: '' },
      rules: {
        companyName: [{ required: true, message: '请选择企业' }],
        companyCode: [{ required: true, message: '请输入信用代码' }],
        usePurpose: [{ required: true, message: '请选择使用目的' }],
        operator: [{ required: true, message: '请输入经办人' }]
      },
      pendingCurrentPage: 1,
      pendingPageSize: 10,
      usedCurrentPage: 1,
      usedPageSize: 10
    }
  },
  computed: {
    pendingList() {
      return this.allCerts.filter(cert => {
        const govt = cert.factory_input || {}
        const ent = cert.driver_input || {}
        return govt.fac_contactNumber === '通过' && !ent.dr_txid
      })
    },
    usedList() {
      return this.allCerts.filter(cert => {
        const ent = cert.driver_input || {}
        return ent.dr_txid && ent.dr_txid !== ''
      })
    },
    pagedPendingList() {
      const start = (this.pendingCurrentPage - 1) * this.pendingPageSize
      return this.pendingList.slice(start, start + this.pendingPageSize)
    },
    pagedUsedList() {
      const start = (this.usedCurrentPage - 1) * this.usedPageSize
      return this.usedList.slice(start, start + this.usedPageSize)
    }
  },
  mounted() {
    this.loadCerts()
  },
  methods: {
    loadCerts() {
      this.loading = true
      getAllFruitInfo().then(res => {
        if (res.code === 200 && res.data) {
          this.allCerts = JSON.parse(res.data)
        }
      }).finally(() => { this.loading = false })
    },
    openUse(cert) {
      this.currentCert = cert
      this.use = { companyName: '', companyCode: '', usePurpose: '', operator: '' }
      this.dialogVisible = true
    },
    doUse() {
      this.$refs.useForm.validate(valid => {
        if (!valid) return
        const fd = new FormData()
        fd.append('traceability_code', this.currentCert.traceability_code)
        fd.append('companyName', this.use.companyName)
        fd.append('companyCode', this.use.companyCode)
        fd.append('usePurpose', this.use.usePurpose)
        fd.append('operator', this.use.operator)
        request({ url: '/enterpriseUse', method: 'post', data: fd }).then(res => {
          if (res.code === 200) {
            this.$message.success('备案成功！')
            this.dialogVisible = false
            this.loadCerts()
          } else {
            this.$message.error(res.message || '备案失败')
          }
        })
      })
    },
    viewDetail(row) {
      this.$router.push({ path: '/trace', query: { code: row.traceability_code }})
    },
    handlePendingPageChange(page) {
      this.pendingCurrentPage = page
    },
    handleUsedPageChange(page) {
      this.usedCurrentPage = page
    }
  }
}
</script>

<style lang="scss" scoped>
.enterprise-container {
  padding: 20px;
  background: #f0f2f5;
  min-height: 100vh;
  width: 100%;
  box-sizing: border-box;
}
h3 {
  text-align: center;
  margin-bottom: 10px;
}
.role-info {
  margin-bottom: 20px;
  font-size: 14px;
  color: #606266;
  text-align: center;
}
.stat-cards {
  margin-bottom: 20px;
}
.list-card {
  border-radius: 16px;
  overflow: hidden;
  width: 100%;
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
.el-table {
  width: 100%;
}
</style>
