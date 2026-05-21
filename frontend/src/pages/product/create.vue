<script setup lang="ts">
import { ref } from "vue"
import { createProduct } from "@/api/product"

const title = ref("")
const description = ref("")
const price = ref("")
const originalPrice = ref("")
const category = ref("教材")
const condition = ref("轻微使用")
const contact = ref("")
const images = ref<string[]>([])
const submitting = ref(false)

const categories = ["教材", "电子", "生活", "运动", "其他"]
const conditions = ["全新", "几乎全新", "轻微使用", "明显使用痕迹"]

async function handleSubmit() {
  if (!title.value.trim()) { uni.showToast({ title: "请输入标题", icon: "none" }); return }
  if (!price.value || parseFloat(price.value) <= 0) { uni.showToast({ title: "请输入有效价格", icon: "none" }); return }

  submitting.value = true
  try {
    const data: any = {
      title: title.value.trim(),
      description: description.value.trim(),
      price: parseFloat(price.value),
      category: category.value,
      condition: condition.value,
      contact: contact.value.trim(),
      images: images.value,
    }
    if (originalPrice.value) data.original_price = parseFloat(originalPrice.value)

    await createProduct(data)
    uni.showToast({ title: "发布成功", icon: "success" })
    setTimeout(() => uni.navigateBack(), 1000)
  } catch {
    uni.showToast({ title: "发布失败", icon: "none" })
  }
  submitting.value = false
}

function chooseImage() {
  uni.chooseImage({
    count: 9,
    success(res) {
      images.value = images.value.concat(res.tempFilePaths).slice(0, 9)
    },
  })
}

function removeImage(idx: number) {
  images.value.splice(idx, 1)
}
</script>

<template>
  <view class="container">
    <view class="form-group">
      <text class="form-label">商品标题 <text class="required">*</text></text>
      <input v-model="title" class="form-input" placeholder="请输入商品标题" maxlength="50" />
    </view>

    <view class="form-group">
      <text class="form-label">描述</text>
      <textarea v-model="description" class="form-textarea" placeholder="描述商品详情、使用情况等" maxlength="500" />
    </view>

    <view class="form-row">
      <view class="form-group flex-1">
        <text class="form-label">价格 (元) <text class="required">*</text></text>
        <input v-model="price" class="form-input" placeholder="0.00" type="digit" />
      </view>
      <view class="form-group flex-1">
        <text class="form-label">原价 (可选)</text>
        <input v-model="originalPrice" class="form-input" placeholder="0.00" type="digit" />
      </view>
    </view>

    <view class="form-row">
      <view class="form-group flex-1">
        <text class="form-label">分类</text>
        <picker :value="categories.indexOf(category)" :range="categories" @change="(e: any) => category = categories[e.detail.value]">
          <view class="picker">{{ category }}</view>
        </picker>
      </view>
      <view class="form-group flex-1">
        <text class="form-label">成色</text>
        <picker :value="conditions.indexOf(condition)" :range="conditions" @change="(e: any) => condition = conditions[e.detail.value]">
          <view class="picker">{{ condition }}</view>
        </picker>
      </view>
    </view>

    <view class="form-group">
      <text class="form-label">联系方式</text>
      <input v-model="contact" class="form-input" placeholder="微信/手机号/QQ" />
    </view>

    <view class="form-group">
      <text class="form-label">图片 (最多9张)</text>
      <view class="image-list">
        <view v-for="(img, i) in images" :key="i" class="image-item">
          <image :src="img" class="preview-img" mode="aspectFill" />
          <text class="remove-img" @click="removeImage(i)">✕</text>
        </view>
        <view v-if="images.length < 9" class="image-item add-img" @click="chooseImage">
          <text class="add-icon">+</text>
        </view>
      </view>
    </view>

    <button class="submit-btn" :loading="submitting" :disabled="submitting" @click="handleSubmit">
      发布
    </button>
  </view>
</template>

<style lang="scss">
page { background: #f5f7fa; }
.container { padding: 30rpx; }

.form-group { margin-bottom: 24rpx; }
.form-label {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 12rpx;
  display: block;
}
.required { color: #e74c3c; }
.form-input {
  background: #fff;
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  font-size: 28rpx;
  height: 44rpx;
}
.form-textarea {
  background: #fff;
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  font-size: 28rpx;
  width: 100%;
  min-height: 160rpx;
  box-sizing: border-box;
}
.form-row { display: flex; gap: 20rpx; }
.flex-1 { flex: 1; }
.picker {
  background: #fff;
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  font-size: 28rpx;
  color: #333;
}

.image-list { display: flex; flex-wrap: wrap; gap: 16rpx; }
.image-item {
  width: 180rpx;
  height: 180rpx;
  border-radius: 12rpx;
  overflow: hidden;
  position: relative;
}
.preview-img { width: 100%; height: 100%; }
.remove-img {
  position: absolute;
  top: 4rpx;
  right: 4rpx;
  width: 36rpx;
  height: 36rpx;
  background: rgba(0,0,0,0.5);
  color: #fff;
  text-align: center;
  line-height: 36rpx;
  border-radius: 50%;
  font-size: 20rpx;
}
.add-img {
  background: #fff;
  border: 2rpx dashed #ccc;
  display: flex;
  align-items: center;
  justify-content: center;
}
.add-icon { font-size: 48rpx; color: #ccc; }

.submit-btn {
  width: 100%;
  background: linear-gradient(135deg, #C67A6A, #1E2A3A);
  color: #fff;
  border: none;
  border-radius: 50rpx;
  font-size: 32rpx;
  height: 88rpx;
  line-height: 88rpx;
  margin-top: 40rpx;
}
</style>
