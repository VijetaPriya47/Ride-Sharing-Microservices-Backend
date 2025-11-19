import { PaymentEventSessionCreatedData } from "../contracts"
import { Button } from "./ui/button"
import { loadStripe } from "@stripe/stripe-js"

interface StripePaymentButtonProps {
  paymentSession: PaymentEventSessionCreatedData
  isLoading?: boolean
}

// Initialize Stripe
const stripePromise = loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY!)

export const StripePaymentButton = ({
  paymentSession,
  isLoading = false,
}: StripePaymentButtonProps) => {
  const handlePayment = async () => {
    // Check if we are in mock mode
    if (paymentSession.sessionID.startsWith("cs_test_mock_session_")) {
      console.log("Mock session detected, simulating successful payment...");
      window.location.href = "/?payment=success";
      return;
    }

    const stripe = await stripePromise

    if (!stripe) {
      console.error("Stripe failed to load")
      return
    }

    // Redirect to Stripe Checkout
    const { error } = await stripe.redirectToCheckout({ sessionId: paymentSession.sessionID })

    if (error) {
      console.error("Payment error:", error)
    }
  }

  if (!process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY) {
    return (
      <Button
        disabled
        className="w-full bg-red-500 text-white"
      >
        Stripe API KEY is not set on the NEXTJS app
      </Button>
    )
  }

  return (
    <Button
      onClick={handlePayment}
      disabled={isLoading}
      className="w-full"
    >
      {isLoading ? "Loading..." : `Pay ${paymentSession.amount} ${paymentSession.currency}`}
    </Button>
  )
} 