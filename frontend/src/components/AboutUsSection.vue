<template>
    <section class="about-section" id="about-us">
        <div class="about-window">
            <div class="about-titlebar">
                <div class="about-titlebar-left">
                    <img
                        src="/favicon.png"
                        alt=""
                        class="about-icon"
                        width="22"
                        height="22"
                    >
                    <span class="about-title">About Us</span>
                </div>
            </div>
            <div class="about-body">
                <p>
                    ChiGame.Site is a <strong>collection of games</strong> you can play right in your browser —
                    <strong>no install, no sign-up</strong>, just pick a game and start playing.
                </p>
                <p>
                    I built this site with an
                    <strong>offline-first</strong> approach so the games stay fast and reliable. More titles are on
                    the way.
                </p>
            </div>
        </div>
    </section>
</template>

<style scoped>
.about-section {
    display: flex;
    justify-content: center;
    margin-top: 4rem;
}

.about-window {
    width: min(100%, 36rem);
    border: 1px solid #0054e3;
    border-radius: 10px 10px 12px 12px;
    overflow: hidden;
    box-shadow:
        1px 1px 0 #000,
        0 8px 32px rgba(0, 20, 60, 0.35),
        0 2px 8px rgba(0, 0, 0, 0.2);
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
}

.about-titlebar {
    display: flex;
    align-items: center;
    padding: 6px 8px;
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
}

.about-titlebar-left {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 0;
}

.about-icon {
    width: 22px;
    height: 22px;
    flex-shrink: 0;
    filter: drop-shadow(1px 1px 1px rgba(0, 0, 0, 0.45));
}

.about-title {
    font-size: 0.95rem;
    font-weight: bold;
    color: #fff;
    text-shadow: 1px 1px #000;
    line-height: 1;
}

.about-body {
    padding: 1.125rem 1.25rem 1.25rem;
    background: #ece9d8;
    border-top: 1px solid #0054e3;
    box-shadow: inset 1px 1px 0 #fff;
}

.about-body p {
    margin: 0;
    font-size: clamp(0.9rem, 2.5vw, 1rem);
    line-height: 1.6;
    font-weight: 400;
    color: #000;
}

.about-body p + p {
    margin-top: 0.75rem;
}

.about-body strong {
    font-weight: 700;
}
</style>
