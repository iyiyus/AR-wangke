<template>
  <div class="landing">
    <!-- 背景光斑 -->
    <div class="bg-orb orb1"></div>
    <div class="bg-orb orb2"></div>
    <div class="bg-orb orb3"></div>

    <!-- 顶部导航 -->
    <header class="landing-header">
      <div class="header-inner">
        <div class="brand">
          <img v-if="siteConfig.logo" :src="siteConfig.logo" class="brand-logo" />
          <span class="brand-name">{{ siteConfig.sitename || '网课代刷管理系统' }}</span>
        </div>
        <div class="header-nav">
          <a href="#features" @click.prevent="scrollTo('features')">功能介绍</a>
          <a href="#advantages" @click.prevent="scrollTo('advantages')">平台优势</a>
          <a href="#" @click.prevent="goLogin">登录</a>
        </div>
      </div>
    </header>

    <!-- Hero -->
    <section class="hero">
      <div class="hero-inner">
        <div class="hero-badge">✨ 稳定运营 · 全年无休</div>
        <h1 class="hero-title">
          一站式<span class="grad-text">网课代刷</span><br />聚合下单平台
        </h1>
        <p class="hero-desc">支持多平台网课代刷、实习打卡、盖章病历聚合管理。<br />极速处理，稳定可靠，让代刷业务更简单。</p>
        <div class="hero-btns">
          <button class="btn-primary" @click="goLogin">
            立即进入系统
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none"><path d="M5 12h14M13 6l6 6-6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </button>
          <button class="btn-ghost" @click="scrollTo('features')">了解更多</button>
        </div>
        <div class="hero-stats">
          <div class="stat-item" v-for="s in stats" :key="s.label">
            <b>{{ s.value }}</b><span>{{ s.label }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- 功能 -->
    <section id="features" class="features">
      <div class="section-inner">
        <div class="section-head">
          <h2 class="section-title">核心功能</h2>
          <p class="section-desc">全方位满足网课代刷业务需求</p>
        </div>
        <div class="feature-grid">
          <div class="feature-card" v-for="f in features" :key="f.title">
            <div class="feature-icon" :style="{ color: f.color }">
              <span v-html="f.icon"></span>
            </div>
            <h3>{{ f.title }}</h3>
            <p>{{ f.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 优势 -->
    <section id="advantages" class="advantages">
      <div class="section-inner">
        <div class="section-head">
          <h2 class="section-title">为什么选择我们</h2>
          <p class="section-desc">专业团队，多年沉淀</p>
        </div>
        <div class="adv-list">
          <div class="adv-card" v-for="a in advantages" :key="a.title">
            <div class="adv-icon" v-html="a.icon"></div>
            <h4>{{ a.title }}</h4>
            <p>{{ a.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 支持平台 -->
    <section class="platforms">
      <div class="section-inner">
        <div class="section-head">
          <h2 class="section-title">支持平台</h2>
          <p class="section-desc">聚合8200+网课平台，一站搞定</p>
        </div>
        <div class="platform-grid">
          <div class="platform-chip" v-for="p in platforms" :key="p">{{ p }}</div>
        </div>
      </div>
    </section>

    <!-- 使用流程 -->
    <section class="steps">
      <div class="section-inner">
        <div class="section-head">
          <h2 class="section-title">三步开始使用</h2>
          <p class="section-desc">简单快捷，马上上手</p>
        </div>
        <div class="step-list">
          <div class="step-item" v-for="(s, i) in steps" :key="s.title">
            <div class="step-num">{{ i + 1 }}</div>
            <h4>{{ s.title }}</h4>
            <p>{{ s.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="cta">
      <div class="cta-card">
        <h2>开始使用，畅享便捷代刷体验</h2>
        <p>注册即用 · 稳定可靠 · 技术支持全程护航</p>
        <button class="btn-primary btn-lg" @click="goLogin">立即登录</button>
      </div>
    </section>

    <footer class="landing-footer">
      <p>© 2026 {{ siteConfig.sitename || '网课代刷管理系统' }}. All rights reserved.</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
  import { onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { siteConfig, loadSiteConfig } from '@/store/modules/site'

  defineOptions({ name: 'HomePage' })

  const router = useRouter()

  const stats = [
    { value: '8200+', label: '支持平台' },
    { value: '99.9%', label: '成功率' },
    { value: '24h', label: '全天处理' },
    { value: '5000+', label: '代理用户' }
  ]

  const features = [
    { icon: '📚', title: '网课代刷', desc: '支持8200+网课平台聚合下单，一键批量自动刷课', color: '#3b82f6', glow: 'rgba(59,130,246,.12)' },
    { icon: '📝', title: '实习打卡', desc: '工学云、慧职教、习讯云等多平台自动打卡', color: '#10b981', glow: 'rgba(16,185,129,.12)' },
    { icon: '🖋️', title: '盖章服务', desc: '4000+企业公章，线上盖章邮寄到家', color: '#8b5cf6', glow: 'rgba(139,92,246,.12)' },
    { icon: '🏥', title: '病历模板', desc: '标准病历模板在线提交自动生成', color: '#f59e0b', glow: 'rgba(245,158,11,.12)' },
    { icon: '💰', title: '代理分销', desc: '多级代理体系密价管理自动分佣', color: '#ec4899', glow: 'rgba(236,72,153,.12)' },
    { icon: '📊', title: '数据统计', desc: '订单流水收益报表一目了然', color: '#06b6d4', glow: 'rgba(6,182,212,.12)' }
  ]

  const advantages = [
    { icon: '🛡️', title: '稳定可靠', desc: '多年运营经验，多通道备援，确保订单稳定处理' },
    { icon: '⚡', title: '极速处理', desc: '全天24小时极速处理，当天发货当天完成' },
    { icon: '🔒', title: '安全保障', desc: '数据加密传输，订单隐私保护，售后无忧' },
    { icon: '🎧', title: '技术支持', desc: '专业团队在线答疑，快速响应各种问题' }
  ]

  const platforms = ['智慧树', '学习通', '知到', '中国大学MOOC', '雨课堂', '学堂在线', '超星尔雅', 'iSmart', 'U校园', '智慧职教', '工学云', '习讯云', ' more...']

  const steps = [
    { title: '注册账号', desc: '快速注册，立即开通系统权限' },
    { title: '选择服务', desc: '网课代刷/实习打卡/盖章病历，按需下单' },
    { title: '坐等结果', desc: '系统自动处理，完成后自动通知' }
  ]

  onMounted(() => loadSiteConfig())
  const goLogin = () => router.push('/auth/login')
  const scrollTo = (id: string) => document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
</script>

<style scoped>
  .landing {
    position: relative;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    background: linear-gradient(180deg, #e0f2fe 0%, #f0f9ff 30%, #fafbff 60%, #fff 100%);
    overflow-x: hidden;
  }

  /* 背景光斑 */
  .bg-orb { position: fixed; border-radius: 50%; filter: blur(100px); z-index: 0; pointer-events: none; }
  .orb1 { width: 600px; height: 600px; background: rgba(59,130,246,.12); top: -200px; left: -100px; }
  .orb2 { width: 500px; height: 500px; background: rgba(139,92,246,.08); top: 40%; right: -150px; }
  .orb3 { width: 400px; height: 400px; background: rgba(6,182,212,.08); bottom: 10%; left: 30%; }

  /* Header */
  .landing-header {
    position: fixed; top: 0; left: 0; right: 0; z-index: 100;
    background: rgba(255,255,255,.75); backdrop-filter: blur(20px);
    border-bottom: 1px solid rgba(0,0,0,.04);
  }
  .header-inner {
    max-width: 1200px; margin: 0 auto; padding: 0 24px;
    height: 64px; display: flex; align-items: center; justify-content: space-between;
  }
  .brand { display: flex; align-items: center; gap: 10px; }
  .brand-logo { width: 32px; height: 32px; border-radius: 8px; }
  .brand-name { font-size: 18px; font-weight: 700; color: #111827; }
  .header-nav { display: flex; gap: 32px; align-items: center; }
  .header-nav a { color: #4b5563; text-decoration: none; font-size: 14px; font-weight: 500; transition: color .2s; }
  .header-nav a:hover { color: #3b82f6; }

  /* Hero */
  .hero { position: relative; z-index: 1; padding: 140px 0 100px; text-align: center; }
  .hero-inner { max-width: 900px; margin: 0 auto; padding: 0 24px; }
  .hero-badge {
    display: inline-block; padding: 6px 16px; border-radius: 20px;
    background: rgba(59,130,246,.08); color: #3b82f6; font-size: 13px; font-weight: 600;
    margin-bottom: 24px;
  }
  .hero-title { font-size: 56px; font-weight: 800; line-height: 1.15; color: #111827; margin: 0 0 24px; letter-spacing: -1px; }
  .grad-text {
    background: linear-gradient(135deg, #3b82f6, #8b5cf6);
    -webkit-background-clip: text; -webkit-text-fill-color: transparent;
  }
  .hero-desc { font-size: 18px; color: #6b7280; line-height: 1.7; margin: 0 0 40px; }
  .hero-btns { display: flex; gap: 16px; justify-content: center; margin-bottom: 56px; }
  .btn-primary {
    display: inline-flex; align-items: center; gap: 8px;
    padding: 14px 32px; border-radius: 12px; border: none;
    background: linear-gradient(135deg, #2563eb, #3b82f6); color: #fff;
    font-size: 16px; font-weight: 600; cursor: pointer;
    box-shadow: 0 4px 16px rgba(59,130,246,.3);
    transition: transform .2s, box-shadow .2s;
  }
  .btn-primary:hover { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(59,130,246,.4); }
  .btn-lg { padding: 16px 48px; font-size: 18px; }
  .btn-ghost {
    padding: 14px 32px; border-radius: 12px; border: 1.5px solid #e5e7eb;
    background: #fff; color: #374151; font-size: 16px; font-weight: 600; cursor: pointer;
    transition: all .2s;
  }
  .btn-ghost:hover { border-color: #3b82f6; color: #3b82f6; }
  .hero-stats { display: flex; gap: 48px; justify-content: center; }
  .stat-item b { display: block; font-size: 32px; color: #111827; }
  .stat-item span { font-size: 14px; color: #9ca3af; }

  /* Features */
  .features { position: relative; z-index: 1; padding: 80px 0; }
  .section-inner { max-width: 1200px; margin: 0 auto; padding: 0 24px; }
  .section-head { text-align: center; margin-bottom: 56px; }
  .section-title { font-size: 36px; font-weight: 700; color: #111827; margin: 0 0 8px; }
  .section-desc { color: #9ca3af; font-size: 16px; margin: 0; }
  .feature-grid { display: grid; grid-template-columns: repeat(3,1fr); gap: 24px; }
  .feature-card {
    position: relative; padding: 36px 28px; border-radius: 20px;
    background: #fff; border: 1px solid #f0f0f5; overflow: hidden;
    transition: transform .25s, box-shadow .25s;
  }
  .feature-card:hover { transform: translateY(-6px); box-shadow: 0 16px 40px rgba(0,0,0,.08); }
  .feature-glow { position: absolute; top: 0; left: 0; right: 0; height: 4px; }
  .feature-icon { font-size: 28px; margin-bottom: 16px; }
  .feature-card h3 { font-size: 18px; margin: 0 0 8px; color: #111827; }
  .feature-card p { font-size: 14px; color: #6b7280; line-height: 1.7; margin: 0; }

  /* Advantages */
  .advantages { position: relative; z-index: 1; padding: 80px 0; }
  .adv-list { display: grid; grid-template-columns: repeat(2,1fr); gap: 24px; }
  .adv-card {
    padding: 32px; border-radius: 16px; background: #fff;
    border: 1px solid #f0f0f5;
  }
  .adv-icon { font-size: 32px; margin-bottom: 12px; }
  .adv-card h4 { font-size: 18px; margin: 0 0 8px; color: #111827; }
  .adv-card p { font-size: 14px; color: #6b7280; margin: 0; line-height: 1.7; }

  /* 支持平台 */
  .platforms { position: relative; z-index: 1; padding: 80px 0; }
  .platform-grid { display: flex; flex-wrap: wrap; gap: 12px; justify-content: center; }
  .platform-chip {
    padding: 10px 20px; border-radius: 24px; background: #fff;
    border: 1px solid #e5e7eb; font-size: 14px; color: #4b5563;
    transition: all .2s;
  }
  .platform-chip:hover { border-color: #3b82f6; color: #3b82f6; }

  /* 使用流程 */
  .steps { position: relative; z-index: 1; padding: 80px 0; }
  .step-list { display: grid; grid-template-columns: repeat(3,1fr); gap: 24px; }
  .step-item { text-align: center; padding: 40px 24px; border-radius: 20px; background: #fff; border: 1px solid #f0f0f5; }
  .step-num {
    width: 56px; height: 56px; margin: 0 auto 20px; border-radius: 50%;
    background: linear-gradient(135deg, #2563eb, #3b82f6); color: #fff;
    font-size: 24px; font-weight: 700; display: flex; align-items: center; justify-content: center;
  }
  .step-item h4 { font-size: 18px; margin: 0 0 8px; color: #111827; }
  .step-item p { font-size: 14px; color: #6b7280; margin: 0; }

  /* CTA */
  .cta { position: relative; z-index: 1; padding: 80px 0; text-align: center; }
  .cta-card {
    max-width: 800px; margin: 0 auto; padding: 60px 40px;
    border-radius: 24px; background: linear-gradient(135deg, #2563eb, #3b82f6);
    box-shadow: 0 20px 60px rgba(59,130,246,.3);
  }
  .cta-card h2 { font-size: 32px; color: #fff; margin: 0 0 12px; }
  .cta-card p { color: rgba(255,255,255,.85); margin: 0 0 32px; }
  .cta .btn-primary { background: #fff; color: #2563eb; box-shadow: none; }

  .landing-footer { position: relative; z-index: 1; padding: 32px 0; text-align: center; color: #9ca3af; font-size: 13px; }

  @media (max-width: 768px) {
    .hero-title { font-size: 30px; }
    .feature-grid { grid-template-columns: 1fr; }
    .adv-list { grid-template-columns: 1fr; }
    .hero-stats { gap: 20px; flex-wrap: wrap; }
    .stat-item b { font-size: 24px; }
  }
</style>
