<template>
    <div class="xp-window">
        <slot></slot>
    </div>
</template>

<style scoped>
.xp-window {
    --title-btn-size: 32px;
    --title-icon-size: 20px;
    --title-font-size: 14px;
    display: flex;
    flex-direction: column;
    align-items: stretch;
    border: 1px solid #0054e3;
    border-radius: 8px 8px 0 0;
    box-shadow: 1px 1px 0 #000;
    overflow: hidden;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    width: fit-content;
    max-width: 100%;
}

@media (max-width: 640px) {
    .xp-window {
        --title-btn-size: 34px;
        --title-icon-size: 22px;
        --title-font-size: 16px;
        width: 100%;
    }
}
</style>
