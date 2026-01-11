<!-- Source: https://docs.wise.com/api-reference -->

Last updated  1 month ago

# Wise Platform API Reference
The Wise Platform API is a REST-based interface that enables programmatic access to Wise's payment infrastructure. All endpoints return JSON-formatted responses and use standard HTTP methods and status codes.

New to wise?

We strongly recommend first reading our **[Getting Started Guide](/guides/developer)** to help you set up credentials and make your first call.

## Before you begin

To use this API reference effectively, you should have:

- Received Valid [API credentials from Wise](/guides/developer/auth-and-security) (Client ID and Client Secret)
- Understand OAuth 2.0 authentication
- Be familiar with RESTful API concepts

## Core API resources

| Resource | Purpose |
| --- | --- |
| Quote | Exchange rate and fee calculations |
| Recipient | Beneficiary account management |
| Transfer | Payment creation and execution |
| Balance | Multi-currency account operations |
| Profile | Account ownership details |

**Not sure which workflow to build?** Start with our [Integration Guides](/guides/product/send-money/use-cases) for step-by-step implementation examples.