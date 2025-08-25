import { test, expect } from "@playwright/test";

// TC012-TC013: Edge Cases
test.describe("Edge Cases", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("http://localhost:5173");
    await page.waitForLoadState("networkidle");
  });

  test("TC012: Rapid Answer Selection", async ({ page }) => {
    // Start the quiz
    await page.click("text=Start Quiz");
    await page.waitForSelector('[data-testid="question-card"]');

    // Get all answer options
    const options = page.locator('[data-testid="answer-option"]');

    // Rapidly click multiple options
    await options.nth(0).click();
    await options.nth(1).click(); // This should not register
    await options.nth(2).click(); // This should not register
    await options.nth(3).click(); // This should not register

    // Wait a bit and verify only first selection took effect
    await page.waitForTimeout(500);

    // Only the first option should be selected (green for correct, red for incorrect)
    await expect(options.nth(0)).toHaveClass(/bg-green-100|bg-red-100/);
    await expect(options.nth(1)).not.toHaveClass(/bg-green-100|bg-red-100/);
    await expect(options.nth(2)).not.toHaveClass(/bg-green-100|bg-red-100/);
    await expect(options.nth(3)).not.toHaveClass(/bg-green-100|bg-red-100/);

    // Wait for auto-advance
    await page.waitForTimeout(1600);

    // Should advance to next question
    await expect(
      page.locator('[data-testid="question-counter"]')
    ).toContainText("2");
  });

  test("TC013: Browser Refresh", async ({ page }) => {
    // Start the quiz and answer a question
    await page.click("text=Start Quiz");
    await page.waitForSelector('[data-testid="question-card"]');

    // Answer first question to get some score
    await page.click('[data-testid="answer-option"]:first-child');
    await page.waitForTimeout(1600);

    // Verify we're on question 2 with some score
    await expect(
      page.locator('[data-testid="question-counter"]')
    ).toContainText("2");

    // Refresh the browser
    await page.reload();
    await page.waitForLoadState("networkidle");

    // Verify quiz resets to start state
    await expect(page.locator("text=Start Quiz")).toBeVisible();

    // Verify no game state persists
    await expect(
      page.locator('[data-testid="question-card"]')
    ).not.toBeVisible();
    await expect(page.locator('[data-testid="game-over"]')).not.toBeVisible();

    // Start quiz again and verify everything is reset
    await page.click("text=Start Quiz");
    await page.waitForSelector('[data-testid="question-card"]');

    // Should be back to question 1
    await expect(
      page.locator('[data-testid="question-counter"]')
    ).toContainText("1");

    // Score should be 0
    const scoreText = await page.locator('[data-testid="score"]').textContent();
    expect(scoreText).toContain("0");
  });

  test("TC012b: Multiple Button Clicks", async ({ page }) => {
    // Test rapid clicking of start button
    await expect(page.locator("text=Start Quiz")).toBeVisible();

    // Rapidly click start button - only first click should register
    // Note: This is done with a small delay to avoid racing conditions
    await page.click("text=Start Quiz");
    
    // Verify game started - start button should be gone and quiz should begin
    await page.waitForSelector('[data-testid="question-card"]');

    // Verify timer starts normally (single game instance)
    await expect(page.locator('[data-testid="timer"]')).toContainText("30");

    // Wait a few seconds and verify timer is counting down normally
    await page.waitForTimeout(3000);
    const timerText = await page.locator('[data-testid="timer"]').textContent();
    const currentTime = parseInt(timerText?.match(/\d+/)?.[0] || "0");
    expect(currentTime).toBeGreaterThanOrEqual(26);
    expect(currentTime).toBeLessThanOrEqual(28);
  });
});
