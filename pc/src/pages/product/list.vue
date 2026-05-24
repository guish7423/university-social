<template>
  <div class="products-page">
    <div class="page-header">
      <h1>二手</h1>
      <el-button type="primary" @click="$router.push('/product/create')">发布闲置</el-button>
    </div>
    <div v-if="loading && !products.length" class="loading-wrap"><el-skeleton :rows="4" animated /></div>
    <div v-else-if="!products.length" class="empty-state"><el-empty description="暂无商品" /></div>
    <template v-else>
      <div class="product-grid">
        <div v-for="p in products" class="product-card stagger-item" :key="p.id" @click="$router.push('/products/' + p.id)">
          <div class="img-wrap">
            <img v-if="p.images?.length" :src="p.images[0]" loading="lazy" />
            <div v-else class="no-img">暂无图片</div>
          </div>
          <div class="card-body">
            <h3>{{ p.title }}</h3>
            <div class="price">
              <span class="current">¥{{ p.price }}</span>
              <span v-if="p.original_price" class="original">¥{{ p.original_price }}</span>
            </div>
            <div class="meta">
              <span>{{ p.category }}</span>
              <span>{{ p.condition }}</span>
            </div>
          </div>
        </div>
      </div>
      <div v-if="hasMore" class="load-more">
        <el-button :loading="loading" @click="loadMore" text>加载更多</el-button>
      </div>
      <div v-if="!hasMore && products.length > 0" class="no-more">没有更多了</div>
</template>

  </div>
</template>


<script setup lang="ts">
import type { ProductData } from "@/api/product"
import { usePagination } from "@/composables/usePagination"
import { listProducts } from "@/api/product"

const { items: products, loading, hasMore, loadMore } = usePagination<ProductData>({
  fetchFn: (offset, limit) => listProducts({ offset, limit }),
})

loadMore()
</script>

<style scoped lang="scss">
@use "@/styles/variables.scss" as *;

.products-page { max-width: 900px; }
.page-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px;
  h1 { font-size: 22px; font-weight: 700; }
}

.product-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(210px, 1fr)); gap: 14px; }

.product-card {
  background: $bg-card; border: 1px solid $border-color;
  border-radius: $radius-md; overflow: hidden; cursor: pointer; transition: all 0.2s;
  &:hover { border-color: $primary-light; transform: translateY(-2px); box-shadow: 0 4px 16px rgba(0,0,0,0.3); }

  .img-wrap {
    height: 160px; overflow: hidden;
    img { width: 100%; height: 100%; object-fit: cover; }
    .no-img { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: $text-muted; font-size: 13px; background: $bg-sidebar; }
  }

  .card-body { padding: 12px;
    h3 { font-size: 14px; font-weight: 600; margin-bottom: 6px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
    .price { margin-bottom: 6px;
      .current { font-size: 16px; font-weight: 700; color: $primary; }
      .original { font-size: 12px; color: $text-muted; text-decoration: line-through; margin-left: 6px; }
    }
    .meta { display: flex; gap: 8px; font-size: 12px; color: $text-muted; }
  }
}
.load-more { text-align: center; padding: 20px 0; }
.no-more { text-align: center; padding: 20px 0; color: $text-muted; font-size: 13px; }
</style>
