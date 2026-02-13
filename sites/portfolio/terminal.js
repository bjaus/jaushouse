document.addEventListener('DOMContentLoaded', () => {
  const terminal = document.getElementById('terminal');
  const sections = Array.from(document.querySelectorAll('.section'));
  const TYPE_SPEED = 40;
  const PAUSE_AFTER_CMD = 300;
  const PAUSE_BETWEEN_SECTIONS = 400;

  // Hide all sections
  sections.forEach(s => s.style.display = 'none');

  // Persistent cursor element that moves with the action
  const cursorLine = document.createElement('div');
  cursorLine.className = 'cmd-line';
  cursorLine.innerHTML = '<span class="prompt">$</span><span class="cursor"></span>';
  terminal.insertBefore(cursorLine, sections[0]);

  function scrollToCursor() {
    cursorLine.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
  }

  function typeText(text) {
    return new Promise(resolve => {
      const typed = document.createElement('span');
      typed.className = 'cmd';
      // Insert typed text before the cursor span
      const cursor = cursorLine.querySelector('.cursor');
      cursorLine.insertBefore(typed, cursor);

      let i = 0;
      const interval = setInterval(() => {
        typed.textContent += text[i];
        i++;
        scrollToCursor();
        if (i >= text.length) {
          clearInterval(interval);
          resolve();
        }
      }, TYPE_SPEED);
    });
  }

  function wait(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  function resetCursorLine() {
    // Clear everything except prompt and cursor
    cursorLine.innerHTML = '<span class="prompt">$</span><span class="cursor"></span>';
  }

  async function runSection(section, isLast) {
    const cmd = section.dataset.cmd;
    const output = section.querySelector('.output');
    output.style.display = 'none';

    // Type the command into the cursor line
    await typeText(cmd);
    await wait(PAUSE_AFTER_CMD);

    // Freeze the typed command line (remove cursor, make it static)
    const frozenLine = document.createElement('div');
    frozenLine.className = 'cmd-line';
    frozenLine.innerHTML = '<span class="prompt">$</span><span class="cmd">' + cmd + '</span>';

    // Show the section with the frozen command and output
    section.style.display = 'block';
    section.insertBefore(frozenLine, output);
    output.style.display = 'block';
    output.classList.add('visible');

    // Move cursor line after this section
    if (!isLast) {
      const spacer = document.createElement('div');
      spacer.className = 'spacer';
      section.after(spacer);
      spacer.after(cursorLine);
    } else {
      section.after(cursorLine);
    }

    // Reset cursor line for next command
    resetCursorLine();
    scrollToCursor();

    if (!isLast) {
      await wait(PAUSE_BETWEEN_SECTIONS);
    }
  }

  async function run() {
    await wait(500);
    for (let i = 0; i < sections.length; i++) {
      await runSection(sections[i], i === sections.length - 1);
    }

    // Stagger fade-in of tech logos after typing completes
    await wait(500);
    const techPanel = document.getElementById('tech-panel');
    if (techPanel) {
      const icons = techPanel.querySelectorAll('img');
      icons.forEach((icon, i) => {
        setTimeout(() => icon.classList.add('visible'), i * 120);
      });
    }
  }

  run();
});
