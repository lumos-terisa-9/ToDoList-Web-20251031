<template>
  <div class="home-page">
    <!-- 广告轮播区域 -->
    <div class="ad-section">
      <div class="ad-carousel">
        <div class="ad-track" :style="{ transform: `translateX(-${currentAdIndex * 100}%)` }">
          <div v-for="(ad, index) in ads" :key="index" class="ad-slide">
            <img :src="ad.image" :alt="ad.title" class="ad-image">
            <div class="ad-content">
              <h3 class="ad-title">{{ ad.title }}</h3>
              <p class="ad-description">{{ ad.description }}</p>
            </div>
          </div>
        </div>

        <!-- 轮播指示器 -->
        <div class="carousel-indicators">
          <button
            v-for="(ad, index) in ads"
            :key="index"
            :class="['indicator', { active: currentAdIndex === index }]"
            @click="currentAdIndex = index"
          ></button>
        </div>

        <!-- 导航按钮 -->
        <button class="carousel-btn prev" @click="prevAd">‹</button>
        <button class="carousel-btn next" @click="nextAd">›</button>
      </div>
    </div>

    <!-- 介绍和问答区域 -->
    <div class="content-section">
      <!-- 网站介绍 -->
      <div class="intro-section">
        <h2 class="section-title">关于我们</h2>
        <div class="intro-content">
          <div class="intro-text">
            <p>欢迎使用 TeamTaskHub - 您的智能团队任务管理助手！我们致力于帮助团队更高效地协作，让项目管理变得简单而愉快。</p>
            <div class="features">
              <div class="feature-item">
                <span class="feature-icon">📅</span>
                <span>智能日历规划</span>
              </div>
              <div class="feature-item">
                <span class="feature-icon">👥</span>
                <span>团队协作管理</span>
              </div>
              <div class="feature-item">
                <span class="feature-icon">⚡</span>
                <span>实时进度跟踪</span>
              </div>
              <div class="feature-item">
                <span class="feature-icon">🔔</span>
                <span>智能提醒通知</span>
              </div>
            </div>
          </div>
          <div class="intro-stats">
            <div class="stat-item">
              <div class="stat-number">10,000+</div>
              <div class="stat-label">活跃团队</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">98%</div>
              <div class="stat-label">用户满意度</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">24/7</div>
              <div class="stat-label">技术支持</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 常见问题 -->
      <div class="faq-section">
        <h2 class="section-title">常见问题</h2>
        <div class="faq-list">
          <div
            v-for="(faq, index) in faqs"
            :key="index"
            class="faq-item"
            :class="{ active: activeFaq === index }"
            @click="toggleFaq(index)"
          >
            <div class="faq-question">
              <span>{{ faq.question }}</span>
              <span class="faq-icon">{{ activeFaq === index ? '−' : '+' }}</span>
            </div>
            <div class="faq-answer" v-if="activeFaq === index">
              <p>{{ faq.answer }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 快速开始 -->
      <div class="quick-start-section">
        <h2 class="section-title">快速开始</h2>
        <div class="quick-steps">
          <div class="step">
            <div class="step-number">1</div>
            <div class="step-content">
              <h3>创建项目</h3>
              <p>点击右上角按钮创建您的第一个项目</p>
            </div>
          </div>
          <div class="step">
            <div class="step-number">2</div>
            <div class="step-content">
              <h3>添加任务</h3>
              <p>在日历视图中添加和管理任务</p>
            </div>
          </div>
          <div class="step">
            <div class="step-number">3</div>
            <div class="step-content">
              <h3>邀请成员</h3>
              <p>邀请团队成员共同协作</p>
            </div>
          </div>
          <div class="step">
            <div class="step-number">4</div>
            <div class="step-content">
              <h3>开始协作</h3>
              <p>享受高效的团队协作体验</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

// 广告数据
const ads = ref([
  {
    image: 'https://images.unsplash.com/photo-1611224923853-80b023f02d71?w=800&h=400&fit=crop',
    title: '智能团队协作',
    description: '让团队协作更高效、更智能'
  },
  {
    image: 'https://images.unsplash.com/photo-1542744173-8e7e53415bb0?w=800&h=400&fit=crop',
    title: '项目管理新体验',
    description: '直观的界面，强大的功能'
  },
  {
    image: 'https://images.unsplash.com/photo-1552664730-d307ca884978?w=800&h=400&fit=crop',
    title: '随时随地办公',
    description: '支持多设备同步，随时随地管理任务'
  }
])

// 常见问题数据
const faqs = ref([
  {
    question: '如何创建新的项目？',
    answer: '点击页面右上角的"新建项目"按钮，填写项目基本信息即可创建新项目。您可以为项目设置名称、描述、开始和结束时间等。'
  },
  {
    question: '支持多少团队成员协作？',
    answer: '我们支持无限数量的团队成员协作。根据您的套餐类型，可以享受不同的存储空间和高级功能。'
  },
  {
    question: '数据安全如何保障？',
    answer: '我们采用银行级别的数据加密技术，所有数据传输都经过SSL加密，确保您的数据安全。'
  },
  {
    question: '是否支持移动端？',
    answer: '是的，我们提供完整的移动端支持，您可以在iOS和Android设备上使用我们的应用。'
  },
  {
    question: '如何导出项目数据？',
    answer: '在项目设置中，您可以找到数据导出功能，支持导出为Excel、PDF等多种格式。'
  }
])

const currentAdIndex = ref(0)
const activeFaq = ref(0)
let autoPlayTimer = null

// 广告轮播控制
function nextAd() {
  currentAdIndex.value = (currentAdIndex.value + 1) % ads.value.length
}

function prevAd() {
  currentAdIndex.value = (currentAdIndex.value - 1 + ads.value.length) % ads.value.length
}

// 自动轮播
function startAutoPlay() {
  autoPlayTimer = setInterval(nextAd, 5000)
}

function stopAutoPlay() {
  if (autoPlayTimer) {
    clearInterval(autoPlayTimer)
    autoPlayTimer = null
  }
}

// FAQ切换
function toggleFaq(index) {
  activeFaq.value = activeFaq.value === index ? -1 : index
}

onMounted(() => {
  startAutoPlay()
})

onUnmounted(() => {
  stopAutoPlay()
})
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  position: relative;
  z-index: 1;
}

/* 广告轮播区域 */
.ad-section {
  padding: 120px 20px 70px 20px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.ad-carousel {
  position: relative;
  max-width: 1300px;
  margin: 0 auto;
  overflow: hidden;
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

.ad-track {
  display: flex;
  transition: transform 0.5s ease;
}

.ad-slide {
  flex: 0 0 100%;
  position: relative;
  height: 350px;
}

.ad-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.ad-content {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 30px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.8));
}

.ad-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

.ad-description {
  font-size: 1.2rem;
  opacity: 0.9;
}

/* 轮播指示器 */
.carousel-indicators {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 10px;
}

.indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 2px solid white;
  background: transparent;
  cursor: pointer;
  transition: all 0.3s ease;
}

.indicator.active {
  background: white;
}

/* 导航按钮 */
.carousel-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 50px;
  height: 50px;
  border: none;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  font-size: 1.5rem;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.carousel-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-50%) scale(1.1);
}

.carousel-btn.prev {
  left: 20px;
}

.carousel-btn.next {
  right: 20px;
}

/* 内容区域 */
.content-section {
  max-width: 1200px;
  margin: 0 auto;
  padding: 60px 20px;
}

.section-title {
  font-size: 2.5rem;
  font-weight: 700;
  text-align: center;
  margin-bottom: 40px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

/* 介绍区域 */
.intro-section {
  margin-bottom: 80px;
}

.intro-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 60px;
  align-items: start;
}

.intro-text p {
  font-size: 1.2rem;
  line-height: 1.8;
  margin-bottom: 30px;
  opacity: 0.9;
}

.features {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  backdrop-filter: blur(10px);
}

.feature-icon {
  font-size: 1.5rem;
}

.intro-stats {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
  backdrop-filter: blur(10px);
}

.stat-number {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 1rem;
  opacity: 0.8;
}

/* FAQ区域 */
.faq-section {
  margin-bottom: 80px;
}

.faq-list {
  max-width: 800px;
  margin: 0 auto;
}

.faq-item {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
  margin-bottom: 15px;
  overflow: hidden;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.faq-item:hover {
  background: rgba(255, 255, 255, 0.15);
}

.faq-question {
  padding: 20px 25px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  font-size: 1.1rem;
  font-weight: 600;
}

.faq-icon {
  font-size: 1.5rem;
  font-weight: 300;
}

.faq-answer {
  padding: 0 25px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
}

.faq-answer p {
  line-height: 1.6;
  opacity: 0.9;
  margin-top: 15px;
}

/* 快速开始区域 */
.quick-steps {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 30px;
}

.step {
  text-align: center;
  padding: 30px 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(10px);
  transition: transform 0.3s ease;
}

.step:hover {
  transform: translateY(-5px);
}

.step-number {
  width: 60px;
  height: 60px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 auto 20px;
}

.step-content h3 {
  font-size: 1.3rem;
  margin-bottom: 10px;
}

.step-content p {
  opacity: 0.8;
  line-height: 1.6;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .ad-title {
    font-size: 2rem;
  }

  .ad-description {
    font-size: 1rem;
  }

  .intro-content {
    grid-template-columns: 1fr;
    gap: 40px;
  }

  .features {
    grid-template-columns: 1fr;
  }

  .quick-steps {
    grid-template-columns: 1fr;
  }

  .carousel-btn {
    width: 40px;
    height: 40px;
    font-size: 1.2rem;
  }
}
</style>
