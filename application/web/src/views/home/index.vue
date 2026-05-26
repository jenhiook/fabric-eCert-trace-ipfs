<template>
  <div class="home-container">
    <!-- 标题 -->
    <div class="hero-section">
      <h1>基于区块链的电子证照存储与溯源系统</h1>
      <p class="hero-desc">Hyperledger Fabric + IPFS 分布式存储 · 全流程可信追溯 · 多角色协同治理</p>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <i class="el-icon-s-grid stat-icon" style="color:#409eff" />
          <div class="stat-num">{{ stats.blocks }}</div>
          <div class="stat-text">区块高度</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <i class="el-icon-s-data stat-icon" style="color:#67c23a" />
          <div class="stat-num">{{ stats.certs }}</div>
          <div class="stat-text">证照总数</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <i class="el-icon-user-solid stat-icon" style="color:#e6a23c" />
          <div class="stat-num">4</div>
          <div class="stat-text">参与角色</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <i class="el-icon-link stat-icon" style="color:#909399" />
          <div class="stat-num">4</div>
          <div class="stat-text">共识节点</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 全生命周期流程 -->
    <el-card style="margin-top:30px">
      <div slot="header">电子证照全生命周期流程</div>
      <div class="flow-line">
        <div class="flow-step" :class="{ active: role === '个人用户' }">
          <div class="flow-icon"><i class="el-icon-edit" /></div>
          <div class="flow-title">证照上链</div>
          <div class="flow-role">个人用户</div>
          <div class="flow-desc">上传证照原件、填写实名信息、生成溯源码，写入区块链</div>
        </div>
        <div class="flow-arrow"><i class="el-icon-right" /></div>
        <div class="flow-step" :class="{ active: role === '政务部门' }">
          <div class="flow-icon"><i class="el-icon-check" /></div>
          <div class="flow-title">部门审核</div>
          <div class="flow-role">政务部门</div>
          <div class="flow-desc">核验证照真实性，审核通过后批准上链，审核信息存证</div>
        </div>
        <div class="flow-arrow"><i class="el-icon-right" /></div>
        <div class="flow-step" :class="{ active: role === '企业组织' }">
          <div class="flow-icon"><i class="el-icon-s-order" /></div>
          <div class="flow-title">使用备案</div>
          <div class="flow-role">企业组织</div>
          <div class="flow-desc">查询审核通过的证照，填写使用目的备案，不可篡改</div>
        </div>
        <div class="flow-arrow"><i class="el-icon-right" /></div>
        <div class="flow-step" :class="{ active: role === '技术支撑实体' }">
          <div class="flow-icon"><i class="el-icon-circle-check" /></div>
          <div class="flow-title">技术核验</div>
          <div class="flow-role">技术支撑实体</div>
          <div class="flow-desc">CA等实体对证照进行密码学核验，验证结果上链</div>
        </div>
      </div>
    </el-card>

    <!-- 图片展示 -->
    <el-row :gutter="20" style="margin-top:30px">
      <el-col :span="12">
        <el-card shadow="hover">
          <div slot="header">系统架构</div>
          <img src="@/assets/arch.png" style="width:100%;border-radius:6px" alt="系统架构图">
          <p class="img-note">▲ 架构示意图</p>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <div slot="header">业务流程</div>
          <img src="@/assets/flow.png" style="width:100%;border-radius:6px" alt="业务流程图">
          <p class="img-note">▲ 流程示意图</p>
        </el-card>
      </el-col>
    </el-row>

    <!-- 技术介绍 -->
    <el-row :gutter="20" style="margin-top:30px">
      <el-col :span="8">
        <el-card shadow="hover" class="feature-card">
          <div slot="header"><i class="el-icon-document" /> 证照全生命周期管理</div>
          <p>个人上传 → 政务审核 → 企业备案 → 技术核验，四阶段全流程可追溯。</p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="feature-card">
          <div slot="header"><i class="el-icon-c-scale-to-original" /> 区块链可信存证</div>
          <p>Hyperledger Fabric 联盟链，元数据上链不可篡改，交易哈希唯一可验证。</p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="feature-card">
          <div slot="header"><i class="el-icon-cloudy" /> IPFS 分布式存储</div>
          <p>证照原图存储于 IPFS 星际文件系统，链上仅存内容指纹（CID）。</p>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷开始 -->
    <el-card style="margin-top:30px">
      <div slot="header">快捷开始</div>
      <div style="text-align:center">
        <el-tag v-if="role" :type="roleTagType" size="medium" effect="dark" style="margin-bottom:12px">
          当前身份：{{ role }}
        </el-tag>
        <p v-if="roleDesc" style="color:#666;margin-bottom:16px;font-size:14px">{{ roleDesc }}</p>

        <el-button v-if="role === '个人用户'" type="primary" icon="el-icon-edit" style="margin:4px" @click="$router.push('/uplink')">证照上链录入</el-button>
        <el-button v-if="role === '政务部门'" type="warning" icon="el-icon-check" style="margin:4px" @click="$router.push('/govt-audit')">证照审核审批</el-button>
        <el-button v-if="role === '企业组织'" type="info" icon="el-icon-s-order" style="margin:4px" @click="$router.push('/enterprise')">证照使用备案</el-button>
        <el-button v-if="role === '技术支撑实体'" type="success" icon="el-icon-circle-check" style="margin:4px" @click="$router.push('/tech-support')">证照技术核验</el-button>

        <el-button type="primary" plain icon="el-icon-search" style="margin:4px" @click="$router.push('/trace')">溯源查询</el-button>
        <el-button type="primary" plain icon="el-icon-data-line" style="margin:4px" @click="$router.push('/dashboard')">可视化分析</el-button>
      </div>
    </el-card>

    <p style="text-align:center;margin-top:30px;color:#999;font-size:13px">
      基于区块链的电子证照存储与溯源系统 © 2026 &nbsp;|&nbsp; 吴晓汀 &nbsp;|&nbsp; 翁茹芳
    </p>
  </div>
</template>

<script>
import { getStats } from '@/api/trace'

export default {
  data() {
    return { stats: { blocks: '--', certs: '--' }}
  },
  computed: {
    role() { return this.$store.state.user.userType || '' },
    roleTagType() {
      const m = { '个人用户': '', '政务部门': 'warning', '企业组织': 'info', '技术支撑实体': 'success' }
      return m[this.role] || ''
    },
    roleDesc() {
      const m = {
        '个人用户': '您可以上传电子证照、填写实名信息并生成溯源码，证照将写入区块链永久存证。',
        '政务部门': '您可以对待审核的证照进行审批，审核通过或驳回的结果将写入区块链不可篡改。',
        '企业组织': '您可以在招聘、招投标等场景中查询并使用已审核通过的证照进行备案。',
        '技术支撑实体': '您作为CA认证中心等实体，可以对已备案的证照进行密码学核验。'
      }
      return m[this.role] || ''
    }
  },
  mounted() {
    getStats().then(res => {
      if (res.code === 200) {
        this.stats.blocks = res.blockHeight || '--'
        try { this.stats.certs = JSON.parse(res.data).totalCerts || 0 } catch (e) {}
      }
    }).catch(() => {
      // ignore stats error
    })
  }
}
</script>

<style scoped>
.home-container { padding: 30px 40px; max-width: 1200px; margin: 0 auto; }
.hero-section { text-align: center; margin-bottom: 30px; }
.hero-section h1 { font-size: 26px; color: #1a1a2e; margin-bottom: 8px; }
.hero-desc { font-size: 15px; color: #888; }
.stat-card { text-align: center; padding: 10px 0; }
.stat-icon { font-size: 32px; margin-bottom: 8px; }
.stat-num { font-size: 28px; font-weight: bold; color: #303133; }
.stat-text { font-size: 13px; color: #999; margin-top: 4px; }
.flow-line { display: flex; align-items: flex-start; justify-content: center; gap: 0; flex-wrap: wrap; }
.flow-step { width: 160px; text-align: center; padding: 16px 10px; border-radius: 10px; background: #f5f7fa; transition: all 0.3s; }
.flow-step.active { background: #e6f7ff; border: 2px solid #409eff; box-shadow: 0 2px 12px rgba(64,158,255,0.2); }
.flow-icon { font-size: 30px; color: #409eff; margin-bottom: 8px; }
.flow-title { font-size: 16px; font-weight: 600; color: #303133; margin-bottom: 4px; }
.flow-role { font-size: 12px; color: #409eff; margin-bottom: 6px; font-weight: 500; }
.flow-desc { font-size: 12px; color: #888; line-height: 1.6; }
.flow-arrow { display: flex; align-items: center; padding: 0 8px; margin-top: 40px; font-size: 22px; color: #c0c4cc; }
.img-note { text-align: center; color: #bbb; font-size: 12px; margin-top: 8px; }
.feature-card p { color: #666; font-size: 14px; line-height: 1.8; }
</style>
