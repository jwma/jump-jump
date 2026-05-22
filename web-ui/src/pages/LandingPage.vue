<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import {
  Link,
  BarChart3,
  QrCode,
  Globe,
  Users,
  Code2,
  ArrowRight,
  Github,
  Server,
} from 'lucide-vue-next'

const { t, tm } = useI18n()
const router = useRouter()

const iconMap: Record<string, typeof Link> = { Link, BarChart3, QrCode, Globe, Users, Code2 }

function scrollToFeatures() {
  document.getElementById('features')?.scrollIntoView({ behavior: 'smooth' })
}

function goToLogin() {
  router.push({ name: 'login' })
}
</script>

<template>
  <div>
    <!-- Hero -->
    <section class="relative overflow-hidden pt-32 pb-20 sm:pt-40 sm:pb-28">
      <div
        class="pointer-events-none absolute inset-0 bg-gradient-to-br from-blue-50 via-white to-indigo-50 dark:from-blue-950/30 dark:via-gray-950 dark:to-indigo-950/20"
      />
      <div
        class="pointer-events-none absolute top-20 left-1/4 h-72 w-72 rounded-full bg-blue-400/10 blur-3xl dark:bg-blue-500/5"
      />
      <div
        class="pointer-events-none absolute bottom-10 right-1/4 h-72 w-72 rounded-full bg-indigo-400/10 blur-3xl dark:bg-indigo-500/5"
      />

      <div class="relative mx-auto max-w-6xl px-4 text-center sm:px-6 lg:px-8">
        <h1
          class="text-4xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-5xl lg:text-6xl"
        >
          {{ t('landing.hero.title') }}
        </h1>
        <p class="mx-auto mt-6 max-w-2xl text-lg text-gray-600 dark:text-gray-400">
          {{ t('landing.hero.subtitle') }}
        </p>
        <div class="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
          <button
            class="inline-flex cursor-pointer items-center gap-2 rounded-xl bg-blue-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-blue-500/25 transition-colors hover:bg-blue-700 hover:shadow-blue-500/40"
            @click="goToLogin"
          >
            {{ t('landing.hero.cta') }}
            <ArrowRight class="h-5 w-5" />
          </button>
          <button
            class="inline-flex cursor-pointer items-center gap-2 rounded-xl border border-gray-300 px-8 py-3.5 text-base font-semibold text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-gray-800"
            @click="scrollToFeatures"
          >
            {{ t('landing.hero.secondary') }}
          </button>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section id="features" class="py-20 sm:py-28">
      <div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
        <div class="text-center">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white sm:text-4xl">
            {{ t('landing.features.title') }}
          </h2>
          <p class="mt-4 text-lg text-gray-600 dark:text-gray-400">
            {{ t('landing.features.subtitle') }}
          </p>
        </div>

        <div class="mt-16 grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
          <div
            v-for="(item, i) in tm('landing.features.items') as unknown as any[]"
            :key="i"
            class="group rounded-2xl border border-gray-200 bg-white p-6 transition-colors hover:border-blue-300 hover:shadow-lg hover:shadow-blue-500/5 dark:border-gray-800 dark:bg-gray-900 dark:hover:border-blue-700"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-600 transition-colors group-hover:bg-blue-100 dark:bg-blue-950 dark:text-blue-400 dark:group-hover:bg-blue-900"
            >
              <component :is="iconMap[item.icon]" class="h-6 w-6" />
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ item.title }}
            </h3>
            <p class="mt-2 text-sm leading-relaxed text-gray-600 dark:text-gray-400">
              {{ item.description }}
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works -->
    <section
      class="border-y border-gray-200 bg-gray-50 py-20 dark:border-gray-800 dark:bg-gray-900 sm:py-28"
    >
      <div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
        <div class="text-center">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white sm:text-4xl">
            {{ t('landing.howItWorks.title') }}
          </h2>
          <p class="mt-4 text-lg text-gray-600 dark:text-gray-400">
            {{ t('landing.howItWorks.subtitle') }}
          </p>
        </div>

        <div class="mt-16 grid gap-8 sm:grid-cols-3">
          <div
            v-for="(step, i) in tm('landing.howItWorks.steps') as unknown as any[]"
            :key="i"
            class="relative text-center"
          >
            <div
              class="mx-auto mb-6 flex h-14 w-14 items-center justify-center rounded-2xl bg-blue-600 text-xl font-bold text-white"
            >
              {{ step.number }}
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ step.title }}
            </h3>
            <p class="mt-2 text-sm leading-relaxed text-gray-600 dark:text-gray-400">
              {{ step.description }}
            </p>
            <!-- Connector line (hidden on last item and mobile) -->
            <div
              v-if="i < 2"
              class="absolute top-7 left-[calc(50%+40px)] hidden h-px w-[calc(100%-80px)] bg-gray-300 dark:bg-gray-700 sm:block"
            />
          </div>
        </div>
      </div>
    </section>

    <!-- Open Source -->
    <section class="py-20 sm:py-28">
      <div class="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
        <div
          class="rounded-3xl border border-gray-200 bg-gradient-to-br from-blue-50 to-indigo-50 p-10 text-center dark:border-gray-800 dark:from-blue-950/30 dark:to-indigo-950/20 sm:p-16"
        >
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white sm:text-4xl">
            {{ t('landing.openSource.title') }}
          </h2>
          <p class="mx-auto mt-4 max-w-2xl text-lg text-gray-600 dark:text-gray-400">
            {{ t('landing.openSource.subtitle') }}
          </p>
          <div class="mt-8 flex flex-col items-center justify-center gap-4 sm:flex-row">
            <a
              :href="t('landing.openSource.githubUrl')"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex cursor-pointer items-center gap-2 rounded-xl border border-gray-300 bg-white px-6 py-3 text-sm font-semibold text-gray-700 shadow-sm transition-colors hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 dark:hover:bg-gray-800"
            >
              <Github class="h-5 w-5" />
              {{ t('landing.openSource.starGithub') }}
            </a>
            <a
              :href="t('landing.openSource.githubUrl')"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex cursor-pointer items-center gap-2 rounded-xl border border-gray-300 bg-white px-6 py-3 text-sm font-semibold text-gray-700 shadow-sm transition-colors hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 dark:hover:bg-gray-800"
            >
              <Server class="h-5 w-5" />
              {{ t('landing.openSource.selfHost') }}
            </a>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section
      class="border-t border-gray-200 bg-gray-50 py-20 dark:border-gray-800 dark:bg-gray-900 sm:py-28"
    >
      <div class="mx-auto max-w-6xl px-4 text-center sm:px-6 lg:px-8">
        <h2 class="text-3xl font-bold text-gray-900 dark:text-white sm:text-4xl">
          {{ t('landing.cta.title') }}
        </h2>
        <p class="mx-auto mt-4 max-w-xl text-lg text-gray-600 dark:text-gray-400">
          {{ t('landing.cta.subtitle') }}
        </p>
        <button
          class="mt-8 inline-flex cursor-pointer items-center gap-2 rounded-xl bg-blue-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-blue-500/25 transition-colors hover:bg-blue-700 hover:shadow-blue-500/40"
          @click="goToLogin"
        >
          {{ t('landing.cta.button') }}
          <ArrowRight class="h-5 w-5" />
        </button>
      </div>
    </section>
  </div>
</template>
