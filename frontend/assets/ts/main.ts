interface TransliterationRequest {
  text: string;
  input_script?: string;
  output_script: string;
  locale?: string;
}

interface TransliterationResponse {
  id: string;
  input_text: string;
  output_text: string;
  input_script: string;
  output_script: string;
  confidence_score?: number;
}

interface FeedbackRequest {
  suggested_output: string;
  feedback_type: string;
  user_context?: string;
}

class TransliterationService {
  private baseUrl: string;

  constructor(baseUrl: string = '') {
    this.baseUrl = baseUrl;
  }

  async transliterate(request: TransliterationRequest): Promise<TransliterationResponse> {
    const response = await fetch(`${this.baseUrl}/transliterate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  }

  async getTransliteration(id: string): Promise<TransliterationResponse> {
    const response = await fetch(`${this.baseUrl}/transliterate/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  }

  async submitFeedback(id: string, feedback: FeedbackRequest): Promise<void> {
    const response = await fetch(`${this.baseUrl}/transliterate/${id}/feedback`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(feedback),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
  }
}

class TransliterationApp {
  private service: TransliterationService;
  private currentResult: TransliterationResponse | null = null;

  constructor() {
    // Use staging API endpoint
    this.service = new TransliterationService('https://staging-transliterate-5dsi.encr.app');
    this.initializeEventListeners();
  }

  private initializeEventListeners(): void {
    const form = document.getElementById('transliteration-form') as HTMLFormElement;
    const feedbackForm = document.getElementById('feedback-form') as HTMLFormElement;

    if (form) {
      form.addEventListener('submit', this.handleTransliteration.bind(this));
    }

    if (feedbackForm) {
      feedbackForm.addEventListener('submit', this.handleFeedback.bind(this));
    }
  }

  private async handleTransliteration(event: Event): Promise<void> {
    event.preventDefault();
    
    const form = event.target as HTMLFormElement;
    const formData = new FormData(form);
    
    const request: TransliterationRequest = {
      text: formData.get('text') as string,
      output_script: formData.get('output_script') as string,
    };

    const inputScript = formData.get('input_script') as string;
    if (inputScript && inputScript !== 'auto') {
      request.input_script = inputScript;
    }

    try {
      this.showLoading(true);
      this.currentResult = await this.service.transliterate(request);
      this.displayResult(this.currentResult);
    } catch (error) {
      this.displayError(error instanceof Error ? error.message : 'An error occurred');
    } finally {
      this.showLoading(false);
    }
  }

  private async handleFeedback(event: Event): Promise<void> {
    event.preventDefault();

    if (!this.currentResult) {
      this.displayError('No transliteration result to provide feedback for');
      return;
    }

    const form = event.target as HTMLFormElement;
    const formData = new FormData(form);
    
    const feedback: FeedbackRequest = {
      suggested_output: formData.get('suggested_output') as string,
      feedback_type: formData.get('feedback_type') as string,
      user_context: formData.get('user_context') as string || undefined,
    };

    try {
      await this.service.submitFeedback(this.currentResult.id, feedback);
      this.displaySuccessMessage('Feedback submitted successfully!');
      form.reset();
    } catch (error) {
      this.displayError(error instanceof Error ? error.message : 'Failed to submit feedback');
    }
  }

  private displayResult(result: TransliterationResponse): void {
    const resultDiv = document.getElementById('result');
    const feedbackSection = document.getElementById('feedback-section');
    
    if (resultDiv) {
      resultDiv.innerHTML = `
        <h3>Transliteration Result</h3>
        <div class="result-display">
          <div class="user-result">
            <h4>Your Result</h4>
            <div class="result-item">
              <strong>Input:</strong> ${this.escapeHtml(result.input_text)}
            </div>
            <div class="result-item">
              <strong>Output:</strong> <span class="output-text">${this.escapeHtml(result.output_text)}</span>
            </div>
            <div class="result-item">
              <strong>Detected Script:</strong> ${this.escapeHtml(result.input_script)}
            </div>
            ${result.confidence_score ? `<div class="result-item">
              <strong>Confidence:</strong> <span class="confidence">${(result.confidence_score * 100).toFixed(1)}%</span>
            </div>` : ''}
          </div>
          
          <div class="json-result">
            <h4>JSON Response (for judging)</h4>
            <pre><code>${JSON.stringify(result, null, 2)}</code></pre>
          </div>
        </div>
      `;
      resultDiv.style.display = 'block';
    }

    if (feedbackSection) {
      feedbackSection.style.display = 'block';
      const suggestedOutputField = document.getElementById('suggested_output') as HTMLInputElement;
      if (suggestedOutputField) {
        suggestedOutputField.value = result.output_text;
      }
    }
  }

  private displayError(message: string): void {
    const errorDiv = document.getElementById('error');
    if (errorDiv) {
      errorDiv.textContent = message;
      errorDiv.style.display = 'block';
    }
    
    // Hide success message if shown
    const successDiv = document.getElementById('success');
    if (successDiv) {
      successDiv.style.display = 'none';
    }
  }

  private displaySuccessMessage(message: string): void {
    const successDiv = document.getElementById('success');
    if (successDiv) {
      successDiv.textContent = message;
      successDiv.style.display = 'block';
    }
    
    // Hide error message if shown
    const errorDiv = document.getElementById('error');
    if (errorDiv) {
      errorDiv.style.display = 'none';
    }

    // Auto-hide success message after 3 seconds
    setTimeout(() => {
      if (successDiv) {
        successDiv.style.display = 'none';
      }
    }, 3000);
  }

  private showLoading(show: boolean): void {
    const loadingDiv = document.getElementById('loading');
    const submitButton = document.getElementById('submit-button') as HTMLButtonElement;
    
    if (loadingDiv) {
      loadingDiv.style.display = show ? 'block' : 'none';
    }
    
    if (submitButton) {
      submitButton.disabled = show;
      submitButton.textContent = show ? 'Transliterating...' : 'Transliterate';
    }

    // Hide previous results and messages
    if (show) {
      const resultDiv = document.getElementById('result');
      const errorDiv = document.getElementById('error');
      const successDiv = document.getElementById('success');
      
      if (resultDiv) resultDiv.style.display = 'none';
      if (errorDiv) errorDiv.style.display = 'none';
      if (successDiv) successDiv.style.display = 'none';
    }
  }

  private escapeHtml(text: string): string {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }
}

// Initialize the app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  new TransliterationApp();
});