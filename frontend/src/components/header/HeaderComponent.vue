<script setup lang="ts">
import { RouterLink } from 'vue-router';
import { useLocalePath } from '@/composables/useLocalePath';

const { localePath } = useLocalePath();
</script>

<template>
    <header class="site-header mb-8">
        <RouterLink :to="localePath('/')" class="site-hero">
            <div class="site-hero-titlebar">
                <div class="site-hero-titlebar-left">
                    <img
                        src="/favicon.png"
                        alt=""
                        class="site-icon"
                        width="28"
                        height="28"
                    >
                    <span class="site-title">
                        ChiGame<span class="site-tld">.Site</span>
                    </span>
                </div>
                <div class="site-hero-chrome" aria-hidden="true">
                    <span class="chrome-btn chrome-btn-min" />
                    <span class="chrome-btn chrome-btn-max" />
                    <span class="chrome-btn chrome-btn-close" />
                </div>
            </div>
        </RouterLink>
    </header>
</template>

<style scoped>
.site-header {
    display: flex;
    justify-content: center;
    padding: 0 0.75rem;
}

.site-hero {
    display: block;
    width: min(100%, 26rem);
    text-decoration: none;
    color: inherit;
    border: 1px solid #0054e3;
    border-radius: 10px 10px 12px 12px;
    overflow: hidden;
    box-shadow:
        1px 1px 0 #000,
        0 8px 32px rgba(0, 20, 60, 0.35),
        0 2px 8px rgba(0, 0, 0, 0.2);
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.site-hero:hover {
    transform: translateY(-2px);
    box-shadow:
        1px 1px 0 #000,
        0 12px 40px rgba(0, 40, 120, 0.4),
        0 4px 12px rgba(0, 0, 0, 0.25);
}

.site-hero:hover .site-hero-titlebar {
    filter: brightness(1.06);
}

.site-hero:active {
    transform: translateY(0) scale(0.99);
}

.site-hero-titlebar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    padding: 6px 6px 6px 8px;
    background: linear-gradient(
        180deg,
        #0997ff 0%,
        #0053ee 8%,
        #0050ee 40%,
        #06f 88%,
        #06f 93%,
        #005eff 95%,
        #003ddb 96%,
        #003ddb 100%
    );
    transition: filter 0.18s ease;
}

.site-hero-titlebar-left {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 0;
}

.site-icon {
    width: 28px;
    height: 28px;
    flex-shrink: 0;
    filter: drop-shadow(1px 1px 1px rgba(0, 0, 0, 0.45));
}

.site-title {
    font-size: clamp(1.15rem, 4vw, 1.35rem);
    font-weight: bold;
    color: #fff;
    text-shadow: 1px 1px #000;
    line-height: 1;
    white-space: nowrap;
}

.site-tld {
    color: #c8e6ff;
    font-weight: 600;
}

.site-hero-chrome {
    display: flex;
    align-items: center;
    gap: 2px;
    flex-shrink: 0;
}

.chrome-btn {
    display: block;
    width: 21px;
    height: 21px;
    border-radius: 3px;
    position: relative;
    box-sizing: border-box;
}

.chrome-btn-min,
.chrome-btn-max {
    background: linear-gradient(180deg, #3c8cfd 0%, #3b8bfc 50%, #3a8afb 51%, #3586f8 100%);
    border: 1px solid #003c74;
    box-shadow: inset 1px 1px 0 #6ba8ff;
}

.chrome-btn-min::after {
    content: '';
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 55%;
    height: 2px;
    background: #fff;
}

.chrome-btn-max::after {
    content: '';
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 42%;
    height: 42%;
    border: 2px solid #fff;
    box-sizing: border-box;
}

.chrome-btn-close {
    background: linear-gradient(180deg, #e88a8f 0%, #e4787d 50%, #e36b70 51%, #d94a4f 100%);
    border: 1px solid #8b0000;
    box-shadow: inset 1px 1px 0 #f0a8ab;
}

.chrome-btn-close::before,
.chrome-btn-close::after {
    content: '';
    position: absolute;
    left: 50%;
    top: 50%;
    width: 55%;
    height: 2px;
    background: #fff;
}

.chrome-btn-close::before {
    transform: translate(-50%, -50%) rotate(45deg);
}

.chrome-btn-close::after {
    transform: translate(-50%, -50%) rotate(-45deg);
}

.site-hero-body {
    padding: 1rem 1.25rem 1.125rem;
    text-align: center;
    background:
        linear-gradient(
            180deg,
            rgba(255, 255, 255, 0.22) 0%,
            rgba(210, 228, 248, 0.14) 50%,
            rgba(180, 210, 240, 0.1) 100%
        );
    backdrop-filter: blur(14px);
    -webkit-backdrop-filter: blur(14px);
    border-top: 1px solid rgba(255, 255, 255, 0.35);
}

.site-tagline {
    margin: 0;
    font-size: clamp(0.78rem, 2.5vw, 0.88rem);
    font-weight: 600;
    color: #fff;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    text-shadow:
        0 1px 2px rgba(0, 0, 0, 0.75),
        0 0 12px rgba(0, 40, 100, 0.5);
}
</style>
