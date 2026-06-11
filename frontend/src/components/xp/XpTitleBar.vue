<script setup lang="ts">
defineProps<{
    title: string;
    icon: string;
}>();

const emit = defineEmits<{
    close: [];
}>();
</script>

<template>
    <div class="xp-titlebar">
        <div class="xp-titlebar-left">
            <img :src="icon" :alt="`${title} Logo`" class="xp-titlebar-icon">
            <span class="xp-titlebar-text">{{ title }}</span>
        </div>
        <div class="xp-titlebar-buttons">
            <button type="button" class="xp-win-btn xp-win-btn-min" aria-label="Minimize" />
            <button type="button" class="xp-win-btn xp-win-btn-max" aria-label="Maximize" />
            <button
                type="button"
                class="xp-win-btn xp-win-btn-close"
                aria-label="Close"
                @click="emit('close')"
            />
        </div>
    </div>
</template>

<style scoped>
.xp-titlebar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    min-height: calc(var(--title-btn-size) + 8px);
    padding: 4px 4px 4px 5px;
    background: linear-gradient(180deg, #0997ff 0%, #0053ee 8%, #0050ee 40%, #06f 88%, #06f 93%, #005eff 95%, #003ddb 96%, #003ddb 100%);
    border-radius: 8px 8px 0 0;
}

.xp-titlebar-left {
    display: flex;
    align-items: center;
    gap: 4px;
    min-width: 0;
}

.xp-titlebar-icon {
    width: var(--title-icon-size);
    height: var(--title-icon-size);
    flex-shrink: 0;
}

.xp-titlebar-text {
    color: #fff;
    font-size: var(--title-font-size);
    font-weight: bold;
    text-shadow: 1px 1px #000;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1;
}

.xp-titlebar-buttons {
    display: flex;
    align-items: center;
    gap: 2px;
    flex-shrink: 0;
}

.xp-win-btn {
    appearance: none;
    -webkit-appearance: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    position: relative;
    width: var(--title-btn-size);
    height: var(--title-btn-size);
    min-width: var(--title-btn-size);
    max-width: var(--title-btn-size);
    box-sizing: border-box;
    flex-shrink: 0;
    padding: 0;
    margin: 0;
    cursor: default;
    border-radius: 3px;
}

.xp-win-btn-min,
.xp-win-btn-max {
    background: linear-gradient(180deg, #3c8cfd 0%, #3b8bfc 50%, #3a8afb 51%, #3586f8 100%);
    border: 1px solid #003c74;
    box-shadow: inset 1px 1px 0 #6ba8ff;
}

.xp-win-btn-min::after {
    content: '';
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 55%;
    height: 2px;
    background: #fff;
}

.xp-win-btn-max::after {
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

.xp-win-btn-close {
    background: linear-gradient(180deg, #e88a8f 0%, #e4787d 50%, #e36b70 51%, #d94a4f 100%);
    border: 1px solid #8b0000;
    box-shadow: inset 1px 1px 0 #f0a8ab;
}

.xp-win-btn-close::before,
.xp-win-btn-close::after {
    content: '';
    position: absolute;
    left: 50%;
    top: 50%;
    width: 55%;
    height: 2px;
    background: #fff;
}

.xp-win-btn-close::before {
    transform: translate(-50%, -50%) rotate(45deg);
}

.xp-win-btn-close::after {
    transform: translate(-50%, -50%) rotate(-45deg);
}
</style>
